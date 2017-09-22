package transaction

import (
	"encoding/json"
	"sort"

	"github.com/trusch/passchain/state"
	"golang.org/x/crypto/sha3"
)

type SecretAddData struct {
	Secret *state.Secret
}

func (data *SecretAddData) Hash() []byte {
	hash := sha3.New512()
	encoder := json.NewEncoder(hash)
	encoder.Encode(data.Secret.ID)
	encoder.Encode(data.Secret.Value)
	hash.Write(hashShares(data.Secret.Shares))
	return hash.Sum(nil)
}

func hashShares(m map[string]string) []byte {
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
