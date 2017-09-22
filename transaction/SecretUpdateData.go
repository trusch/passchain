package transaction

import (
	"encoding/json"

	"github.com/trusch/passchain/state"
	"golang.org/x/crypto/sha3"
)

type SecretUpdateData struct {
	Secret   *state.Secret
	SenderID string
}

func (data *SecretUpdateData) Hash() []byte {
	hash := sha3.New512()
	encoder := json.NewEncoder(hash)
	encoder.Encode(data.SenderID)
	encoder.Encode(data.Secret.ID)
	encoder.Encode(data.Secret.Value)
	hash.Write(hashShares(data.Secret.Shares))
	return hash.Sum(nil)
}
