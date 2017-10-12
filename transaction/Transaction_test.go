package transaction_test

import (
	"github.com/trusch/passchain/crypto"
	"github.com/trusch/passchain/state"
	. "github.com/trusch/passchain/transaction"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Transaction", func() {
	It("should be possible to find a PoW", func() {
		t := New(AccountAdd, &AccountAddData{Account: &state.Account{}})
		k, _ := crypto.CreateKeyPair()
		Expect(t.Sign(k)).To(Succeed())
		Expect(t.ProofOfWork(16)).To(Succeed())
		Expect(t.Verify(k)).To(Succeed())
		Expect(t.VerifyProofOfWork(16)).To(Succeed())
	})
})
