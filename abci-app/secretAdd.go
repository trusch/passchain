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
	if len(data.Secret.Owners) == 0 {
		return errors.New("no owners supplied")
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
