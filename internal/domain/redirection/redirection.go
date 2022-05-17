package redirection

import (
	"github.com/emanuelefalzone/bitly/internal"
)

type Redirection struct {
	Key      string `validate:"required"`
	Location string `validate:"required,url"`
}

func New(key string, location string) (Redirection, error) {
	value := Redirection{
		Key:      key,
		Location: location,
	}
	err := internal.Validate(value)
	if err != nil {
		return Redirection{}, &internal.Error{Op: "NewRedirection", Err: err}
	}
	return value, nil
}
