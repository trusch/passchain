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
	"encoding/json"
	"log"

	"github.com/tendermint/abci/types"
	"github.com/tendermint/merkleeyes/iavl"
	cmn "github.com/tendermint/tmlibs/common"
	"github.com/trusch/passchain/state"
	"github.com/trusch/passchain/transaction"
)

type Application struct {
	types.BaseApplication

	state *state.State
}

func NewApplication() *Application {
	tree := iavl.NewIAVLTree(0, nil)
	return &Application{state: state.NewStateFromTree(tree)}
}

func (app *Application) Info() (resInfo types.ResponseInfo) {
	return types.ResponseInfo{Data: cmn.Fmt("{\"size\":%v}", app.state.Tree.Size())}
}

func (app *Application) DeliverTx(txBytes []byte) types.Result {
	tx := &transaction.Transaction{}
	if err := tx.FromBytes(txBytes); err != nil {
		return types.ErrUnknownRequest
	}
	switch tx.Type {
	case transaction.AccountAdd:
		{
			if err := deliverAccountAddTransaction(tx, app.state); err != nil {
				return types.Result{Code: types.CodeType_BaseInvalidInput, Log: err.Error()}
			}
		}

	case transaction.AccountDel:
		{
			if err := deliverAccountDelTransaction(tx, app.state); err != nil {
				return types.Result{Code: types.CodeType_BaseInvalidInput, Log: err.Error()}
			}
		}
	case transaction.ReputationGive:
		{
			if err := deliverReputationGiveTransaction(tx, app.state); err != nil {
				return types.Result{Code: types.CodeType_BaseInvalidInput, Log: err.Error()}
			}
		}
	case transaction.SecretAdd:
		{
			if err := deliverSecretAddTransaction(tx, app.state); err != nil {
				return types.Result{Code: types.CodeType_BaseInvalidInput, Log: err.Error()}
			}
		}
	case transaction.SecretUpdate:
		{
			if err := deliverSecretUpdateTransaction(tx, app.state); err != nil {
				return types.Result{Code: types.CodeType_BaseInvalidInput, Log: err.Error()}
			}
		}
	case transaction.SecretDel:
		{
			if err := deliverSecretDelTransaction(tx, app.state); err != nil {
				return types.Result{Code: types.CodeType_BaseInvalidInput, Log: err.Error()}
			}
		}
	case transaction.SecretShare:
		{
			if err := deliverSecretShareTransaction(tx, app.state); err != nil {
				return types.Result{Code: types.CodeType_BaseInvalidInput, Log: err.Error()}
			}
		}
	default:
		{
			return types.Result{Code: types.CodeType_BaseInvalidInput, Log: "unknown transaction type"}
		}
	}
	return types.OK
}

func (app *Application) CheckTx(txBytes []byte) types.Result {
	tx := &transaction.Transaction{}
	if err := tx.FromBytes(txBytes); err != nil {
		return types.ErrUnknownRequest
	}
	switch tx.Type {
	case transaction.AccountAdd:
		{
			if err := checkAccountAddTransaction(tx, app.state); err != nil {
				return types.Result{Code: types.CodeType_BaseInvalidInput, Log: err.Error()}
			}
		}

	case transaction.AccountDel:
		{
			if err := checkAccountDelTransaction(tx, app.state); err != nil {
				return types.Result{Code: types.CodeType_BaseInvalidInput, Log: err.Error()}
			}
		}
	case transaction.ReputationGive:
		{
			if err := checkReputationGiveTransaction(tx, app.state); err != nil {
				return types.Result{Code: types.CodeType_BaseInvalidInput, Log: err.Error()}
			}
		}
	case transaction.SecretAdd:
		{
			if err := checkSecretAddTransaction(tx, app.state); err != nil {
				return types.Result{Code: types.CodeType_BaseInvalidInput, Log: err.Error()}
			}
		}
	case transaction.SecretUpdate:
		{
			if err := checkSecretUpdateTransaction(tx, app.state); err != nil {
				return types.Result{Code: types.CodeType_BaseInvalidInput, Log: err.Error()}
			}
		}
	case transaction.SecretDel:
		{
			if err := checkSecretDelTransaction(tx, app.state); err != nil {
				return types.Result{Code: types.CodeType_BaseInvalidInput, Log: err.Error()}
			}
		}
	case transaction.SecretShare:
		{
			if err := checkSecretShareTransaction(tx, app.state); err != nil {
				return types.Result{Code: types.CodeType_BaseInvalidInput, Log: err.Error()}
			}
		}
	default:
		{
			return types.Result{Code: types.CodeType_BaseInvalidInput, Log: "unknown transaction type"}
		}
	}
	return types.OK
}

func (app *Application) Commit() types.Result {
	hash := app.state.Tree.Hash()
	return types.NewResultOK(hash, "")
}

func (app *Application) Query(reqQuery types.RequestQuery) (resQuery types.ResponseQuery) {
	log.Print("query")
	switch reqQuery.Path {
	case "/account":
		{
			var (
				result interface{}
				err    error
			)
			if reqQuery.Data == nil {
				result, err = app.state.ListAccounts()
				log.Printf("got account list: %+v", result)
			} else {
				result, err = app.state.GetAccount(string(reqQuery.Data))
				log.Printf("got account: %+v", result)
			}
			if err != nil {
				resQuery.Code = types.CodeType_BaseInvalidInput
				resQuery.Log = err.Error()
				return
			}
			bs, _ := json.Marshal(result)
			resQuery.Value = bs
		}
	case "/secret":
		{
			var (
				result interface{}
				err    error
			)
			if reqQuery.Data == nil {
				result, err = app.state.ListSecrets()
				log.Printf("got secret list: %+v", result)
			} else {
				result, err = app.state.GetSecret(string(reqQuery.Data))
				log.Printf("got secret: %+v", result)
			}
			if err != nil {
				resQuery.Code = types.CodeType_BaseInvalidInput
				resQuery.Log = err.Error()
				return
			}
			bs, _ := json.Marshal(result)
			resQuery.Value = bs
		}
	default:
		{
			resQuery.Code = types.CodeType_BaseInvalidInput
			resQuery.Log = "wrong path"
			return
		}
	}
	return
}
