package app

import (
	"errors"
	"log"

	"github.com/mitchellh/mapstructure"
	"github.com/trusch/passchain/state"
	"github.com/trusch/passchain/transaction"
)

func checkSecretUpdateTransaction(tx *transaction.Transaction, state *state.State) error {
	data := &transaction.SecretUpdateData{}
	if err := mapstructure.Decode(tx.Data, data); err != nil {
		return err
	}
	tx.Data = data
	k, err := state.GetAccountPubKey(data.SenderID)
	if err != nil {
		return errors.New("pubkey can't be loaded: " + err.Error())
	}
	log.Printf("verify signature %v of %v", tx.Signature, data.SenderID)
	log.Printf("data: %+v (%T)", data, data)
	if err = tx.Verify(k); err != nil {
		return errors.New("tx can't be verified: " + err.Error())
	}
	secret, err := state.GetSecret(data.Secret.ID)
	if err != nil {
		return err
	}
	if _, ok := secret.Shares[data.SenderID]; !ok {
		return errors.New("sender has no share on this secret")
	}
	return nil
}

func deliverSecretUpdateTransaction(tx *transaction.Transaction, state *state.State) error {
	data := &transaction.SecretUpdateData{}
	if err := mapstructure.Decode(tx.Data, &data); err != nil {
		return err
	}
	return state.SetSecret(data.Secret)
}
