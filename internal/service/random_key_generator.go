package service

import (
	"math/rand"
)

const (
	// Only use lower case char (minus i,j and l) and digits so that the key is easy to dictate
	alphabet = "abcdefghkmnopqrstuvwxyz0123456789"
	// With a 6 chars key we have 1291467969 possible keys
	// This should be enough, be can be increased if requested
	length = 6
)

// Inspired by https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
type RandomKeyGenerator struct {
	rand *rand.Rand
}

func NewRandomKeyGenerator(seed int64) KeyGenerator {
	rand := rand.New(rand.NewSource(seed))
	return RandomKeyGenerator{rand: rand}
}

// The Next Key method simply ignores the location parameter and create a fixed length random key
func (s RandomKeyGenerator) NextKey(location string) string {
	b := make([]byte, length)
	// extract random chars from alphabet
	for i := range b {
		b[i] = alphabet[s.rand.Intn(len(alphabet))]
	}
	return string(b)
}
