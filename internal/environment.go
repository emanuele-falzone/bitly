package internal

import (
	"fmt"
	"os"
)

func GetEnv(key string) (string, error) {
	value, present := os.LookupEnv(key)
	if !present {
		return "", &Error{Code: ErrNotFound, Message: fmt.Sprintf("Environment variable %s is not set!", key)}
	}
	if value == "" {
		return "", &Error{Code: ErrInvalid, Message: fmt.Sprintf("Environment variable %s is empty!", key)}
	}
	return value, nil
}
