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

package crypto

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
	"math/big"
	"strings"

	"github.com/trusch/storage/filter/encryption/ecdhe"
)

type Key struct {
	pub  *ecdsa.PublicKey
	priv *ecdsa.PrivateKey
}

func CreateKeyPair() (*Key, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	return &Key{priv: priv, pub: &priv.PublicKey}, nil
}

func NewFromStrings(pub, priv string) (*Key, error) {
	k := &Key{}
	if pub == "" && priv == "" {
		return nil, errors.New("no key material supplied")
	}
	if pub != "" {
		if err := k.SetPubString(pub); err != nil {
			return nil, err
		}
	}
	if priv != "" && pub == "" {
		return nil, errors.New("no pubkey to privkey supplied")
	}
	if priv != "" {
		if err := k.SetPrivString(priv); err != nil {
			return nil, err
		}
	}
	return k, nil
}

func (k *Key) GetPubString() string {
	pub := elliptic.Marshal(elliptic.P256(), k.pub.X, k.pub.Y)
	return base64.StdEncoding.EncodeToString(pub)
}

func (k *Key) GetPrivString() string {
	priv := k.priv.D.Bytes()
	return base64.StdEncoding.EncodeToString(priv)
}

func (k *Key) SetPubString(pub string) error {
	bs, err := base64.StdEncoding.DecodeString(pub)
	if err != nil {
		return err
	}
	x, y := elliptic.Unmarshal(elliptic.P256(), bs)
	k.pub = &ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}
	return nil
}

func (k *Key) SetPrivString(priv string) error {
	bs, err := base64.StdEncoding.DecodeString(priv)
	if err != nil {
		return err
	}
	d := &big.Int{}
	d.SetBytes(bs)
	k.priv = &ecdsa.PrivateKey{D: d, PublicKey: *k.pub}
	k.pub = &k.priv.PublicKey
	return nil
}

func (k *Key) GetWriter(base io.Writer) (io.WriteCloser, error) {
	return ecdhe.NewWriter(base, k.pub)
}

func (k *Key) GetReader(base io.Reader) (io.ReadCloser, error) {
	return ecdhe.NewReader(base, k.priv)
}

func (k *Key) Sign(hash []byte) (string, error) {
	r, s, err := ecdsa.Sign(rand.Reader, k.priv, hash)
	if err != nil {
		return "", err
	}
	rStr := base64.StdEncoding.EncodeToString(r.Bytes())
	sStr := base64.StdEncoding.EncodeToString(s.Bytes())
	return rStr + ":" + sStr, nil
}

func (k *Key) Verify(hash []byte, signature string) error {
	parts := strings.Split(signature, ":")
	if len(parts) != 2 {
		return errors.New("malformed signature")
	}
	rBytes, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return err
	}
	sBytes, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return err
	}
	r := (&big.Int{}).SetBytes(rBytes)
	s := (&big.Int{}).SetBytes(sBytes)
	if !ecdsa.Verify(k.pub, hash, r, s) {
		return errors.New("bad signature")
	}
	return nil
}

func (k *Key) EncryptToString(data []byte) (string, error) {
	buf := &bytes.Buffer{}
	w, err := k.GetWriter(buf)
	if err != nil {
		return "", err
	}
	_, err = w.Write(data)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func (k *Key) DecryptString(cipherText string) ([]byte, error) {
	d, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(d)
	r, err := k.GetReader(buf)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(r)
}
