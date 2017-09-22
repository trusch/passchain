package state

import (
	"encoding/json"
	"errors"

	"github.com/tendermint/tmlibs/merkle"
	"github.com/trusch/passchain/crypto"
)

const (
	accountPrefix = "account::"
	secretPrefix  = "secret::"
)

type State struct {
	Tree merkle.Tree
}

func NewStateFromTree(tree merkle.Tree) *State {
	return &State{tree}
}

type Account struct {
	ID         string         `json:"id" mapstructure:"id"`
	PubKey     string         `json:"pubkey" mapstructure:"pubkey"`
	Reputation map[string]int `json:"reputation" mapstructure:"reputation"`
}

type Secret struct {
	ID     string `json:"id" mapstructure:"id"`
	Value  string
	Shares map[string]string
}

func (s *State) AddAccount(account *Account) error {
	if s.HasAccount(account.ID) {
		return errors.New("account already exists")
	}
	return s.SetAccount(account)
}

func (s *State) SetAccount(account *Account) error {
	bs, err := json.Marshal(account)
	if err != nil {
		return err
	}
	s.Tree.Set([]byte(accountPrefix+account.ID), bs)
	return nil
}

func (s *State) HasAccount(id string) bool {
	return s.Tree.Has([]byte(accountPrefix + id))
}

func (s *State) GetAccount(id string) (*Account, error) {
	_, bs, exists := s.Tree.Get([]byte(accountPrefix + id))
	if !exists {
		return nil, errors.New("no such account")
	}
	acc := &Account{Reputation: make(map[string]int)}
	return acc, json.Unmarshal(bs, acc)
}

func (s *State) DeleteAccount(id string) error {
	_, removed := s.Tree.Remove([]byte(accountPrefix + id))
	if !removed {
		return errors.New("no such account")
	}
	return nil
}

func (s *State) GetAccountPubKey(id string) (*crypto.Key, error) {
	acc, err := s.GetAccount(id)
	if err != nil {
		return nil, err
	}
	return crypto.NewFromStrings(acc.PubKey, "")
}

func (s *State) AddSecret(secret *Secret) error {
	if s.HasSecret(secret.ID) {
		return errors.New("secret already exists")
	}
	return s.SetSecret(secret)
}

func (s *State) SetSecret(secret *Secret) error {
	bs, err := json.Marshal(secret)
	if err != nil {
		return err
	}
	s.Tree.Set([]byte(secretPrefix+secret.ID), bs)
	return nil
}

func (s *State) HasSecret(id string) bool {
	return s.Tree.Has([]byte(secretPrefix + id))
}

func (s *State) GetSecret(id string) (*Secret, error) {
	_, bs, exists := s.Tree.Get([]byte(secretPrefix + id))
	if !exists {
		return nil, errors.New("no such secret")
	}
	acc := &Secret{Shares: make(map[string]string)}
	return acc, json.Unmarshal(bs, acc)
}

func (s *State) DeleteSecret(id string) error {
	_, removed := s.Tree.Remove([]byte(secretPrefix + id))
	if !removed {
		return errors.New("no such secret")
	}
	return nil
}
