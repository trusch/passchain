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
