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

package state

import (
	"encoding/json"
	"errors"

	"github.com/trusch/passchain/crypto"
)

type Account struct {
	ID         string         `json:"id" mapstructure:"id"`
	PubKey     string         `json:"pubkey" mapstructure:"pubkey"`
	Reputation map[string]int `json:"reputation" mapstructure:"reputation"`
}

func (s *State) AddAccount(account *Account) error {
	if s.HasAccount(account.ID) {
		return errors.New("account already exists")
	}
	return s.SetAccount(account)
}

func (s *State) SetAccount(account *Account) error {
	bs, err := json.Marshal(account)
	if err != nil {
		return err
	}
	s.Tree.Set([]byte(accountPrefix+account.ID), bs)
	return nil
}

func (s *State) HasAccount(id string) bool {
	return s.Tree.Has([]byte(accountPrefix + id))
}

func (s *State) GetAccount(id string) (*Account, error) {
	_, bs, exists := s.Tree.Get([]byte(accountPrefix + id))
	if !exists {
		return nil, errors.New("no such account")
	}
	acc := &Account{Reputation: make(map[string]int)}
	return acc, json.Unmarshal(bs, acc)
}

func (s *State) DeleteAccount(id string) error {
	_, removed := s.Tree.Remove([]byte(accountPrefix + id))
	if !removed {
		return errors.New("no such account")
	}
	return nil
}

func (s *State) GetAccountPubKey(id string) (*crypto.Key, error) {
	acc, err := s.GetAccount(id)
	if err != nil {
		return nil, err
	}
	return crypto.NewFromStrings(acc.PubKey, "")
}

func (s *State) ListAccounts() (result []*Account, err error) {
	start := accountPrefix
	end := start[:len(start)-1]
	end = end + string(start[len(start)-1]+1)
	result = make([]*Account, 0)
	s.Tree.IterateRange([]byte(start), []byte(end), true, func(key []byte, value []byte) bool {
		acc := &Account{}
		err = json.Unmarshal(value, acc)
		if err != nil {
			return true
		}
		result = append(result, acc)
		return false
	})
	return
}
