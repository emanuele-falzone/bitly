package redirection

type Repository interface {
	Create(Redirection) error
	Delete(Redirection) error
	FindByKey(string) (Redirection, error)
}
