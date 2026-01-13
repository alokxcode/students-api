package storage

type DB interface {
	Create(name string, email string, password string) (int, error)
}
