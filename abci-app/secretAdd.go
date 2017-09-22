package app

import (
	"errors"

	"github.com/mitchellh/mapstructure"
	"github.com/trusch/passchain/state"
	"github.com/trusch/passchain/transaction"
)

func checkSecretAddTransaction(tx *transaction.Transaction, state *state.State) error {
	data := &transaction.SecretAddData{}
	if err := mapstructure.Decode(tx.Data, data); err != nil {
		return err
	}
	if state.HasSecret(data.Secret.ID) {
		return errors.New("secret exists")
	}
	if len(data.Secret.Shares) == 0 {
		return errors.New("no shares supplied")
	}
	return nil
}

func deliverSecretAddTransaction(tx *transaction.Transaction, state *state.State) error {
	data := &transaction.SecretAddData{}
	if err := mapstructure.Decode(tx.Data, &data); err != nil {
		return err
	}
	return state.AddSecret(data.Secret)
}
