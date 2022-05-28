package service

// KeyGenerator allows for generating new keys
type KeyGenerator interface {
	// NextKey generates a new key
	NextKey() string
}
