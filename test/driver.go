package test

import "context"

// The Driver interface allow us to test the code from different entrypoint
// We just have to adhere to such interface.
type Driver interface {
	CreateRedirection(ctx context.Context, location string) (string, error)
	DeleteRedirection(ctx context.Context, key string) error
	GetRedirectionLocation(ctx context.Context, key string) (string, error)
	GetRedirectionCount(ctx context.Context, key string) (int, error)
	GetRedirectionList(ctx context.Context) ([]string, error)
}
