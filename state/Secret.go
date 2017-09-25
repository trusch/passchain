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
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
)

func (s *State) AddSecret(secret *Secret) error {
	if s.HasSecret(secret.ID) {
		return errors.New("secret already exists")
	}
	return s.SetSecret(secret)
}

func (s *State) SetSecret(secret *Secret) error {
	bs, err := json.Marshal(secret)
	if err != nil {
		return err
	}
	s.Tree.Set([]byte(secretPrefix+secret.ID), bs)
	return nil
}

func (s *State) HasSecret(id string) bool {
	return s.Tree.Has([]byte(secretPrefix + id))
}

func (s *State) GetSecret(id string) (*Secret, error) {
	_, bs, exists := s.Tree.Get([]byte(secretPrefix + id))
	if !exists {
		return nil, errors.New("no such secret")
	}
	acc := &Secret{Shares: make(map[string]string)}
	return acc, json.Unmarshal(bs, acc)
}

func (s *State) DeleteSecret(id string) error {
	_, removed := s.Tree.Remove([]byte(secretPrefix + id))
	if !removed {
		return errors.New("no such secret")
	}
	return nil
}

func (s *State) ListSecrets() (result []*Secret, err error) {
	start := secretPrefix
	end := start[:len(start)-1]
	end = end + string(start[len(start)-1]+1)
	result = make([]*Secret, 0)
	s.Tree.IterateRange([]byte(start), []byte(end), true, func(key []byte, value []byte) bool {
		acc := &Secret{}
		err = json.Unmarshal(value, acc)
		if err != nil {
			return true
		}
		result = append(result, acc)
		return false
	})
	return
}

func (secret *Secret) Encrypt() (aesKey []byte, err error) {
	k := make([]byte, 32)
	if _, err = io.ReadFull(rand.Reader, k); err != nil {
		return nil, err
	}
	key := sha256.Sum256(k)
	return key[:], secret.EncryptWithKey(key[:])
}

func (secret *Secret) EncryptWithKey(key []byte) error {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return err
	}
	iv := make([]byte, aes.BlockSize)
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return err
	}
	buf := &bytes.Buffer{}
	_, err = buf.Write(iv)
	if err != nil {
		return err
	}
	stream := cipher.NewOFB(block, iv[:])
	writer := &cipher.StreamWriter{S: stream, W: buf}
	_, err = writer.Write([]byte(secret.Value))
	if err != nil {
		return err
	}
	secret.Value = base64.StdEncoding.EncodeToString(buf.Bytes())
	return nil
}

func (secret *Secret) Decrypt(key []byte) error {
	valueBytes, err := base64.StdEncoding.DecodeString(secret.Value)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(valueBytes)
	iv := make([]byte, aes.BlockSize)
	bs, err := buf.Read(iv[:])
	if bs != aes.BlockSize {
		return errors.New("ciphertext to short")
	}
	if err != nil {
		return err
	}
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return err
	}
	stream := cipher.NewOFB(block, iv[:])
	reader := &cipher.StreamReader{S: stream, R: buf}
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	secret.Value = string(data)
	return nil
}
