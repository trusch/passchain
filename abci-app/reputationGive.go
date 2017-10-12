/*
 * Copyright (C) 2017 Tino Rusch
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of
 * this software and associated documentation files (the "Software"), to deal in
 * the Software without restriction, including without limitation the rights to
 * use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
 * the Software, and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
 * FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
 * IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

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
	tx.Data = data
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
	if err := tx.VerifyProofOfWork(transaction.DefaultProofOfWorkCost); err != nil {
		return err
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
