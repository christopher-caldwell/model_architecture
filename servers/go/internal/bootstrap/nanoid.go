package bootstrap

import (
	"crypto/rand"
	"math/big"
)

const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_-"

func generateNanoID(size int) string {
	alphabetLen := big.NewInt(int64(len(alphabet)))
	b := make([]byte, size)
	for i := range b {
		n, err := rand.Int(rand.Reader, alphabetLen)
		if err != nil {
			panic("nanoid: rand.Int failed: " + err.Error())
		}
		b[i] = alphabet[n.Int64()]
	}
	return string(b)
}
