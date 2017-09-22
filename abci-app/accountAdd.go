package app

import (
	"errors"

	"github.com/mitchellh/mapstructure"
	"github.com/trusch/passchain/crypto"
	"github.com/trusch/passchain/state"
	"github.com/trusch/passchain/transaction"
)

func checkAccountAddTransaction(tx *transaction.Transaction, state *state.State) error {
	data := &transaction.AccountAddData{}
	if err := mapstructure.Decode(tx.Data, data); err != nil {
		return err
	}
	if state.HasAccount(data.Account.ID) {
		return errors.New("account exists")
	}
	_, err := crypto.NewFromStrings(data.Account.PubKey, "")
	if err != nil {
		return err
	}
	return nil
}

func deliverAccountAddTransaction(tx *transaction.Transaction, state *state.State) error {
	data := &transaction.AccountAddData{}
	if err := mapstructure.Decode(tx.Data, &data); err != nil {
		return err
	}
	return state.AddAccount(data.Account)
}
