package redirection

import (
	"github.com/emanuelefalzone/bitly/internal"
)

// Redirection associates a user provided location to system generated key
type Redirection struct {
	Key      string `validate:"required"`     // The key associated by the system to a given location
	Location string `validate:"required,url"` // The location (URL) that the user want to shorten
}

// New creates a new redirection given the key and location
func New(key, location string) (*Redirection, error) {
	value := &Redirection{
		Key:      key,
		Location: location,
	}

	if err := internal.Validate(value); err != nil {
		return nil, &internal.Error{
			Op:  "Redirection: New",
			Err: err,
		}
	}

	return value, nil
}
