package service

type KeyGenerator interface {
	NextKey(location string) string
}
