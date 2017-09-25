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

package crypto_test

import (
	"bytes"
	"io/ioutil"

	. "github.com/trusch/passchain/crypto"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Key", func() {
	It("should be possible to create a key", func() {
		k, err := CreateKeyPair()
		Expect(err).NotTo(HaveOccurred())
		Expect(k).NotTo(BeNil())
	})

	It("should be possible to serialize + deserialize a key", func() {
		k, err := CreateKeyPair()
		Expect(err).NotTo(HaveOccurred())
		Expect(k).NotTo(BeNil())
		priv := k.GetPrivString()
		pub := k.GetPubString()
		Expect(priv).NotTo(BeEmpty())
		Expect(pub).NotTo(BeEmpty())
		restored, err := NewFromStrings(pub, priv)
		Expect(err).NotTo(HaveOccurred())
		Expect(restored).To(Equal(k))
	})

	It("should be possible to create write only keys (no private key)", func() {
		k, err := CreateKeyPair()
		Expect(err).NotTo(HaveOccurred())
		pub := k.GetPubString()
		restored, err := NewFromStrings(pub, "")
		Expect(err).NotTo(HaveOccurred())
		Expect(restored.GetPubString()).To(Equal(k.GetPubString()))
	})

	It("should be possible to encrypt/decrypt stuff", func() {
		k, err := CreateKeyPair()
		Expect(err).NotTo(HaveOccurred())
		buf := &bytes.Buffer{}
		w, err := k.GetWriter(buf)
		Expect(err).NotTo(HaveOccurred())
		_, err = w.Write([]byte("foobar"))
		Expect(err).NotTo(HaveOccurred())
		err = w.Close()
		Expect(err).NotTo(HaveOccurred())
		reader, err := k.GetReader(buf)
		Expect(err).NotTo(HaveOccurred())
		bs, err := ioutil.ReadAll(reader)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(bs)).To(Equal("foobar"))
	})

	It("should be possible to encrypt/decrypt stuff when the writer only has the pubkey", func() {
		k, err := CreateKeyPair()
		Expect(err).NotTo(HaveOccurred())

		pub := k.GetPubString()
		writeKey, err := NewFromStrings(pub, "")
		Expect(err).NotTo(HaveOccurred())
		buf := &bytes.Buffer{}
		w, err := writeKey.GetWriter(buf)
		Expect(err).NotTo(HaveOccurred())
		_, err = w.Write([]byte("foobar"))
		Expect(err).NotTo(HaveOccurred())
		err = w.Close()
		Expect(err).NotTo(HaveOccurred())

		reader, err := k.GetReader(buf)
		Expect(err).NotTo(HaveOccurred())
		bs, err := ioutil.ReadAll(reader)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(bs)).To(Equal("foobar"))
	})

	It("should be possible to sign/verify stuff", func() {
		k, err := CreateKeyPair()
		Expect(err).NotTo(HaveOccurred())
		hash := []byte("foobar")
		signature, err := k.Sign(hash)
		Expect(err).NotTo(HaveOccurred())

		other, err := NewFromStrings(k.GetPubString(), "")
		Expect(err).NotTo(HaveOccurred())
		err = other.Verify(hash, signature)
		Expect(err).NotTo(HaveOccurred())
	})
})
