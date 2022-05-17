package driver

// The Driver interface allow us to test the code from different entrypoint
// We just have to adhere to such interface
type Driver interface {
	CreateRedirection(location string) (string, error)
	DeleteRedirection(key string) error
	GetRedirectionLocation(key string) (string, error)
}
