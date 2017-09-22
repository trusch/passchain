package app

import (
	"errors"

	"github.com/mitchellh/mapstructure"
	"github.com/trusch/passchain/state"
	"github.com/trusch/passchain/transaction"
)

func checkSecretShareTransaction(tx *transaction.Transaction, state *state.State) error {
	data := &transaction.SecretShareData{}
	if err := mapstructure.Decode(tx.Data, data); err != nil {
		return err
	}
	secret, err := state.GetSecret(data.ID)
	if err != nil {
		return err
	}
	if _, ok := secret.Shares[data.SenderID]; !ok {
		return errors.New("sender has no share on this secret")
	}
	if _, ok := secret.Shares[data.AccountID]; ok {
		return errors.New("share receiver already has a share on this secret")
	}
	k, err := state.GetAccountPubKey(data.SenderID)
	if err != nil {
		return errors.New("pubkey can't be loaded: " + err.Error())
	}
	if err = tx.Verify(k); err != nil {
		return errors.New("tx can't be verified: " + err.Error())
	}
	return nil
}

func deliverSecretShareTransaction(tx *transaction.Transaction, state *state.State) error {
	data := &transaction.SecretShareData{}
	if err := mapstructure.Decode(tx.Data, &data); err != nil {
		return err
	}
	secret, err := state.GetSecret(data.ID)
	if err != nil {
		return err
	}
	secret.Shares[data.AccountID] = data.Key
	return state.SetSecret(secret)
}
