package redirection

import (
	"github.com/emanuelefalzone/bitly/internal"
)

type Redirection struct {
	Key      string `validate:"required"`
	Location string `validate:"required,url"`
}

func New(key, location string) (Redirection, error) {
	value := Redirection{
		Key:      key,
		Location: location,
	}

	if err := internal.Validate(value); err != nil {
		return Redirection{}, &internal.Error{
			Op:  "Redirection: New",
			Err: err,
		}
	}

	return value, nil
}
