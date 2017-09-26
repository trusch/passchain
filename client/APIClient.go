package client

import (
	"errors"
	"fmt"
	"log"

	"github.com/trusch/passchain/crypto"
	"github.com/trusch/passchain/state"
)

// API is the high level interface for passchain client applications
type API interface {
	As(accountID string) (API, error)
	AccountAPI
	ReputationAPI
	SecretAPI
}

// AccountAPI describes all account related functions
type AccountAPI interface {
	CreateAccount(id string) (pub, priv string, err error)
	GetAccount(id string) (*state.Account, error)
	DeleteAccount(id string) error
	ListAccounts(idPrefix string) ([]*state.Account, error)
}

// ReputationAPI describes the reputation related function
type ReputationAPI interface {
	GiveReputation(receiver string, value int) error
}

// SecretAPI describes operations on secrets
type SecretAPI interface {
	CreateSecret(sid string, value string) error
	GetSecret(sid string) (*state.Secret, error)
	DeleteSecret(sid string) error
	ListSecrets(sidPrefix string) ([]*state.Secret, error)
	ShareSecret(sid, accountID string, ownerRights bool) error
	UpdateSecret(sid, value string) error
	UnshareSecret(sid, accountID string) error
}

// NewAPI constructs a new API instances based on an http transport
func NewAPI(endpoint string, key *crypto.Key, account string) API {
	base := NewHTTPClient(endpoint, key, account)
	return &apiClient{endpoint, base}
}

type apiClient struct {
	endpoint string
	base     *BaseClient
}

func (api *apiClient) As(accountID string) (API, error) {
	asAccount, err := api.GetAccount(accountID)
	if err != nil {
		return nil, fmt.Errorf("cannot find account: %v", err)
	}
	privkeySecret, err := api.GetSecret(accountID)
	if err != nil {
		return nil, fmt.Errorf("cannot find shared account secret: %v", err)
	}
	encryptedAESKey, ok := privkeySecret.Shares[api.base.AccountID]
	if !ok {
		return nil, errors.New("no share for us on this secret")
	}
	aesKey, err := api.base.Key.DecryptString(encryptedAESKey)
	if err != nil {
		return nil, err
	}
	err = privkeySecret.Decrypt(aesKey)
	if err != nil {
		return nil, err
	}
	key, err := crypto.NewFromStrings(asAccount.PubKey, string(privkeySecret.Value))
	if err != nil {
		return nil, err
	}
	return NewAPI(api.endpoint, key, accountID), nil
}

func (api *apiClient) CreateAccount(id string) (pub, priv string, err error) {
	key, err := crypto.CreateKeyPair()
	if err != nil {
		return "", "", err
	}
	api.base.Key = key
	if err := api.base.AddAccount(&state.Account{ID: id, PubKey: key.GetPubString()}); err != nil {
		return "", "", err
	}
	return key.GetPubString(), key.GetPrivString(), nil
}

func (api *apiClient) GetAccount(id string) (*state.Account, error) {
	return api.base.GetAccount(id)
}

func (api *apiClient) DeleteAccount(id string) error {
	return api.base.DelAccount(id)
}

func (api *apiClient) ListAccounts(idPrefix string) ([]*state.Account, error) {
	return api.base.ListAccounts()
}

func (api *apiClient) CreateSecret(sid string, value string) error {
	s := &state.Secret{
		ID:     sid,
		Value:  value,
		Shares: make(map[string]string),
		Owners: map[string]bool{
			api.base.AccountID: true,
		},
	}
	aesKey, err := s.Encrypt()
	if err != nil {
		return err
	}
	k := api.base.Key
	encryptedAESKey, err := k.EncryptToString(aesKey)
	if err != nil {
		return err
	}
	s.Shares[api.base.AccountID] = encryptedAESKey
	if err := api.base.AddSecret(s); err != nil {
		return err
	}
	return nil
}

func (api *apiClient) GetSecret(sid string) (*state.Secret, error) {
	secret, err := api.base.GetSecret(sid)
	if err != nil {
		return nil, err
	}
	encryptedAESKey, ok := secret.Shares[api.base.AccountID]
	if ok {
		k := api.base.Key
		aesKey, err := k.DecryptString(encryptedAESKey)
		if err != nil {
			return nil, err
		}
		err = secret.Decrypt(aesKey)
		if err != nil {
			return nil, err
		}
	}
	return secret, nil
}

func (api *apiClient) DeleteSecret(sid string) error {
	return api.base.DelSecret(sid)
}

func (api *apiClient) ListSecrets(sidPrefix string) ([]*state.Secret, error) {
	secrets, err := api.base.ListSecrets()
	if err != nil {
		return nil, err
	}
	for _, s := range secrets {
		if encryptedKey, ok := s.Shares[api.base.AccountID]; ok {
			key, e := api.base.Key.DecryptString(encryptedKey)
			if e != nil {
				log.Fatal(e)
			}
			if e = s.Decrypt(key); e != nil {
				log.Fatal(e)
			}
		}
	}
	return secrets, nil
}

func (api *apiClient) ShareSecret(sid, accountID string, ownerRights bool) error {
	secret, err := api.base.GetSecret(sid)
	if err != nil {
		return err
	}
	encryptedAESKey, ok := secret.Shares[api.base.AccountID]
	if !ok {
		return errors.New("no share for us on this secret")
	}
	k := api.base.Key
	aesKey, err := k.DecryptString(encryptedAESKey)
	if err != nil {
		return err
	}
	acc, err := api.GetAccount(accountID)
	if err != nil {
		log.Print("can not find account " + accountID)
		return err
	}
	otherKey, err := crypto.NewFromStrings(acc.PubKey, "")
	if err != nil {
		return err
	}
	otherEncrptedAESKey, err := otherKey.EncryptToString(aesKey)
	if err != nil {
		return err
	}
	secret.Shares[accountID] = otherEncrptedAESKey
	if ownerRights {
		secret.Owners[accountID] = true
	}
	return api.base.UpdateSecret(secret)
}

func (api *apiClient) UpdateSecret(sid, value string) error {
	sec, err := api.base.GetSecret(sid)
	if err != nil {
		return fmt.Errorf("failed to get secret: %v", err)
	}
	aesKey, err := api.base.Key.DecryptString(sec.Shares[api.base.AccountID])
	if err != nil {
		return fmt.Errorf("failed to decrypt secret: %v", err)
	}
	sec.Value = value
	err = sec.EncryptWithKey(aesKey)
	if err != nil {
		return err
	}
	if err := api.base.UpdateSecret(sec); err != nil {
		return err
	}
	return nil
}

func (api *apiClient) UnshareSecret(sid, accountID string) error {
	secret, err := api.base.GetSecret(sid)
	if err != nil {
		log.Fatal(err)
	}
	delete(secret.Shares, accountID)
	delete(secret.Owners, accountID)
	return api.base.UpdateSecret(secret)
}

func (api *apiClient) GiveReputation(receiver string, value int) error {
	return api.base.GiveReputation(api.base.AccountID, receiver, value)
}
