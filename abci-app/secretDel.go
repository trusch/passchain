package app

import (
	"errors"

	"github.com/mitchellh/mapstructure"
	"github.com/trusch/passchain/state"
	"github.com/trusch/passchain/transaction"
)

func checkSecretDelTransaction(tx *transaction.Transaction, state *state.State) error {
	data := &transaction.SecretDelData{}
	if err := mapstructure.Decode(tx.Data, data); err != nil {
		return err
	}
	if !state.HasSecret(data.ID) {
		return errors.New("secret doesn't exists")
	}
	secret, err := state.GetSecret(data.ID)
	if err != nil {
		return err
	}
	if _, ok := secret.Shares[data.SenderID]; !ok {
		return errors.New("sender has no share on this secret")
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

func deliverSecretDelTransaction(tx *transaction.Transaction, state *state.State) error {
	data := &transaction.SecretDelData{}
	if err := mapstructure.Decode(tx.Data, &data); err != nil {
		return err
	}
	return state.DeleteSecret(data.ID)
}
