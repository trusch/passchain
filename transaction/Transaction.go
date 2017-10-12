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

package transaction

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"sort"
	"time"

	"github.com/trusch/passchain/crypto"

	"golang.org/x/crypto/sha3"
)

type Transaction struct {
	Type      TransactionType `json:"type"`
	Timestamp time.Time       `json:"timestamp"`
	Signature string          `json:"signature"`
	Nonce     uint32          `json:"nonce"`
	Data      interface{}     `json:"data"`
}

type Hashable interface {
	Hash() []byte
}

type TransactionType string

const (
	AccountAdd     TransactionType = "add-account"
	AccountDel     TransactionType = "del-account"
	ReputationGive TransactionType = "give-reputation"
	SecretAdd      TransactionType = "secret-add"
	SecretUpdate   TransactionType = "secret-update"
	SecretDel      TransactionType = "secret-del"
	SecretShare    TransactionType = "secret-share"
)

const DefaultProofOfWorkCost byte = 16

func (t *Transaction) FromBytes(bs []byte) error {
	return json.Unmarshal(bs, t)
}

func (t *Transaction) ToBytes() ([]byte, error) {
	return json.Marshal(t)
}

func (t *Transaction) Hash() []byte {
	hash := sha3.New512()
	encoder := json.NewEncoder(hash)
	encoder.Encode(t.Type)
	encoder.Encode(t.Timestamp)
	if hashable, ok := t.Data.(Hashable); ok {
		hash.Write(hashable.Hash())
	} else {
		encoder.Encode(t.Data)
	}
	return hash.Sum(nil)
}

func (t *Transaction) Sign(key *crypto.Key) error {
	hash := t.Hash()
	signature, err := key.Sign(hash)
	if err != nil {
		return err
	}
	t.Signature = signature
	return nil
}

func (t *Transaction) Verify(key *crypto.Key) error {
	hash := t.Hash()
	return key.Verify(hash, t.Signature)
}

func (t *Transaction) ProofOfWork(cost byte) error {
	for round := 0; round < (1 << 32); round++ {
		t.Nonce = uint32(round)
		if err := t.VerifyProofOfWork(cost); err == nil {
			return nil
		}
	}
	return errors.New("can not find pow")
}

func (t *Transaction) VerifyProofOfWork(cost byte) error {
	hasher := sha3.New512()
	hasher.Write(t.Hash())
	binary.Write(hasher, binary.LittleEndian, t.Nonce)
	tip := uint64(0)
	buf := bytes.NewBuffer(hasher.Sum(nil))
	binary.Read(buf, binary.LittleEndian, &tip)
	if tip<<(64-cost) == 0 {
		return nil
	}
	return errors.New("failed to validate proof of work")
}

func New(t TransactionType, data interface{}) *Transaction {
	return &Transaction{t, time.Now(), "", 0, data}
}

func hashStringMap(m map[string]interface{}) []byte {
	hash := sha3.New512()
	encoder := json.NewEncoder(hash)
	keys := make([]string, len(m))
	i := 0
	for id := range m {
		keys[i] = id
		i++
	}
	sort.Strings(keys)
	for _, key := range keys {
		encoder.Encode(key)
		encoder.Encode(m[key])
	}
	return hash.Sum(nil)
}
