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

package client

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/tendermint/tendermint/rpc/client"
	"github.com/tendermint/tendermint/types"
	"github.com/trusch/passchain/crypto"
	"github.com/trusch/passchain/state"
	"github.com/trusch/passchain/transaction"
)

type Client struct {
	Key       *crypto.Key
	AccountID string
	tm        client.Client
}

func NewHTTPClient(endpoint string, key *crypto.Key, account string) *Client {
	tm := client.NewHTTP(endpoint, "/websocket")
	return &Client{key, account, tm}
}

func (c *Client) AddAccount(acc *state.Account) error {
	tx := transaction.New(transaction.AccountAdd, &transaction.AccountAddData{Account: acc})
	bs, _ := tx.ToBytes()
	res, err := c.tm.BroadcastTxCommit(types.Tx(bs))
	if err != nil {
		return err
	}
	if res.CheckTx.IsErr() {
		return errors.New(res.CheckTx.Error())
	}
	if res.DeliverTx.IsErr() {
		return errors.New(res.DeliverTx.Error())
	}
	return nil
}

func (c *Client) DelAccount(id string) error {
	tx := transaction.New(transaction.AccountDel, &transaction.AccountDelData{ID: id})
	if err := tx.Sign(c.Key); err != nil {
		return err
	}
	bs, _ := tx.ToBytes()
	res, err := c.tm.BroadcastTxCommit(types.Tx(bs))
	if err != nil {
		return err
	}
	if res.CheckTx.IsErr() {
		return errors.New(res.CheckTx.Error())
	}
	if res.DeliverTx.IsErr() {
		return errors.New(res.DeliverTx.Error())
	}
	return nil
}

func (c *Client) GiveReputation(from, to string, value int) error {
	tx := transaction.New(transaction.ReputationGive, &transaction.ReputationGiveData{
		From:  from,
		To:    to,
		Value: value,
	})
	if err := tx.Sign(c.Key); err != nil {
		return err
	}
	bs, _ := tx.ToBytes()
	res, err := c.tm.BroadcastTxCommit(types.Tx(bs))
	if err != nil {
		return err
	}
	if res.CheckTx.IsErr() {
		return errors.New(res.CheckTx.Error())
	}
	if res.DeliverTx.IsErr() {
		return errors.New(res.DeliverTx.Error())
	}
	return nil
}

func (c *Client) GetAccount(id string) (*state.Account, error) {
	resp, err := c.tm.ABCIQuery("/account", []byte(id), false)
	if err != nil {
		return nil, err
	}
	if len(resp.Value) == 0 {
		return nil, errors.New("account not found")
	}
	acc := &state.Account{}
	if err = json.Unmarshal(resp.Value, acc); err != nil {
		return nil, err
	}
	return acc, nil
}

func (c *Client) ListAccounts() ([]*state.Account, error) {
	resp, err := c.tm.ABCIQuery("/account", nil, false)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	if len(resp.Value) == 0 {
		return nil, errors.New("account not found")
	}
	acc := []*state.Account{}
	if err = json.Unmarshal(resp.Value, &acc); err != nil {
		return nil, err
	}
	return acc, nil
}

func (c *Client) ListSecrets() ([]*state.Secret, error) {
	resp, err := c.tm.ABCIQuery("/secret", nil, false)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	if len(resp.Value) == 0 {
		return nil, errors.New("secret not found")
	}
	acc := []*state.Secret{}
	if err = json.Unmarshal(resp.Value, &acc); err != nil {
		return nil, err
	}
	return acc, nil
}

func (c *Client) GetSecret(id string) (*state.Secret, error) {
	resp, err := c.tm.ABCIQuery("/secret", []byte(id), false)
	if err != nil {
		return nil, err
	}
	if len(resp.Value) == 0 {
		return nil, errors.New("secret not found")
	}
	acc := &state.Secret{}
	if err = json.Unmarshal(resp.Value, acc); err != nil {
		return nil, err
	}
	return acc, nil
}

func (c *Client) AddSecret(acc *state.Secret) error {
	tx := transaction.New(transaction.SecretAdd, &transaction.SecretAddData{Secret: acc})
	bs, _ := tx.ToBytes()
	res, err := c.tm.BroadcastTxCommit(types.Tx(bs))
	if err != nil {
		return err
	}
	if res.CheckTx.IsErr() {
		return errors.New(res.CheckTx.Error())
	}
	if res.DeliverTx.IsErr() {
		return errors.New(res.DeliverTx.Error())
	}
	return nil
}

func (c *Client) DelSecret(id string) error {
	tx := transaction.New(transaction.SecretDel, &transaction.SecretDelData{
		ID:       id,
		SenderID: c.AccountID,
	})
	if err := tx.Sign(c.Key); err != nil {
		return err
	}
	bs, _ := tx.ToBytes()
	res, err := c.tm.BroadcastTxCommit(types.Tx(bs))
	if err != nil {
		return err
	}
	if res.CheckTx.IsErr() {
		return errors.New(res.CheckTx.Error())
	}
	if res.DeliverTx.IsErr() {
		return errors.New(res.DeliverTx.Error())
	}
	return nil
}

func (c *Client) UpdateSecret(acc *state.Secret) error {
	tx := transaction.New(transaction.SecretUpdate, &transaction.SecretUpdateData{
		Secret:   acc,
		SenderID: c.AccountID,
	})
	if err := tx.Sign(c.Key); err != nil {
		return err
	}

	bs, _ := tx.ToBytes()
	res, err := c.tm.BroadcastTxCommit(types.Tx(bs))
	if err != nil {
		return err
	}
	if res.CheckTx.IsErr() {
		return errors.New(res.CheckTx.Error())
	}
	if res.DeliverTx.IsErr() {
		return errors.New(res.DeliverTx.Error())
	}
	return nil
}
