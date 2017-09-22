package app

import (
	"errors"

	"github.com/mitchellh/mapstructure"
	"github.com/trusch/passchain/state"
	"github.com/trusch/passchain/transaction"
)

func checkReputationGiveTransaction(tx *transaction.Transaction, state *state.State) error {
	data := &transaction.ReputationGiveData{}
	if err := mapstructure.Decode(tx.Data, data); err != nil {
		return err
	}
	if !state.HasAccount(data.From) {
		return errors.New("reject give-rep because id doesnt exist: " + data.From)
	}
	if !state.HasAccount(data.To) {
		return errors.New("reject give-rep because id doesnt exist: " + data.From)
	}
	if data.Value < -3 || data.Value > 3 {
		return errors.New("reject give-rep because bad value")
	}
	k, err := state.GetAccountPubKey(data.From)
	if err != nil {
		return errors.New("reject give-rep because pubkey cant be loaded")
	}
	if err = tx.Verify(k); err != nil {
		return errors.New("reject give-rep because signature cant be verified")
	}
	return nil
}

func deliverReputationGiveTransaction(tx *transaction.Transaction, state *state.State) error {
	data := &transaction.ReputationGiveData{}
	if err := mapstructure.Decode(tx.Data, &data); err != nil {
		return err
	}
	acc, err := state.GetAccount(data.To)
	if err != nil {
		return err
	}
	if acc.Reputation == nil {
		acc.Reputation = make(map[string]int)
	}
	acc.Reputation[data.From] = data.Value
	err = state.SetAccount(acc)
	if err != nil {
		return err
	}
	return nil
}
