package redirection

import "context"

//go:generate mockgen -destination=../../../test/mock/redirection_repository.go -package=mock -mock_names=Repository=MockRedirectionRepository github.com/emanuelefalzone/bitly/internal/domain/redirection Repository

type Repository interface {
	Create(context.Context, Redirection) error
	Delete(context.Context, Redirection) error
	FindByKey(context.Context, string) (Redirection, error)
	FindAll(context.Context) ([]Redirection, error)
}
