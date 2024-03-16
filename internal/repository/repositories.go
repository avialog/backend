package repository

type Repositories interface {
}

type repositories struct {
}

func NewRepositories() Repositories {
	return &repositories{}
}
