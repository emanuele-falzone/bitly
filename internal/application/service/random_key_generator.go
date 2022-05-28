package service

import (
	"math/rand"
)

const (
	// Only use lower case char (minus i,j and l) and digits so that the key is easy to dictate.
	alphabet = "abcdefghkmnopqrstuvwxyz0123456789"
	// With a 6 chars key we have 1291467969 possible keys.
	length = 6
)

// RandomKeyGenerator creates a fixed length random key
// Inspired by https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
type RandomKeyGenerator struct {
	rand *rand.Rand
}

// NewRandomKeyGenerator creates a new random key generator with the given seed
func NewRandomKeyGenerator(seed int64) *RandomKeyGenerator {
	return &RandomKeyGenerator{rand: rand.New(rand.NewSource(seed))}
}

func (s *RandomKeyGenerator) NextKey() string {
	b := make([]byte, length)
	// extract random chars from alphabet
	for i := range b {
		b[i] = alphabet[s.rand.Intn(len(alphabet))]
	}

	return string(b)
}
