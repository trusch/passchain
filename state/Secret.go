package state

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
)

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
