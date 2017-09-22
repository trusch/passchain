package app

import (
	"errors"

	"github.com/mitchellh/mapstructure"
	"github.com/trusch/passchain/state"
	"github.com/trusch/passchain/transaction"
)

func checkAccountDelTransaction(tx *transaction.Transaction, state *state.State) error {
	data := &transaction.AccountDelData{}
	if err := mapstructure.Decode(tx.Data, data); err != nil {
		return err
	}
	if !state.HasAccount(data.ID) {
		return errors.New("account doesn't exists")
	}
	k, err := state.GetAccountPubKey(data.ID)
	if err != nil {
		return errors.New("pubkey can't be loaded: " + err.Error())
	}
	if err = tx.Verify(k); err != nil {
		return errors.New("tx can't be verified: " + err.Error())
	}
	return nil
}

func deliverAccountDelTransaction(tx *transaction.Transaction, state *state.State) error {
	data := &transaction.AccountDelData{}
	if err := mapstructure.Decode(tx.Data, &data); err != nil {
		return err
	}
	return state.DeleteAccount(data.ID)
}
