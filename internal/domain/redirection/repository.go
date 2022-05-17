package redirection

//go:generate mockgen -destination=../../../test/mock/redirection_repository.go -package=mock -mock_names=Repository=MockRedirectionRepository github.com/emanuelefalzone/bitly/internal/domain/redirection Repository

type Repository interface {
	Create(Redirection) error
	Delete(Redirection) error
	FindByKey(string) (Redirection, error)
}
