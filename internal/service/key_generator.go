package service

//go:generate mockgen -destination=../../test/mock/key_generator_service.go -package=mock github.com/emanuelefalzone/bitly/internal/service KeyGenerator

type KeyGenerator interface {
	NextKey(location string) string
}
