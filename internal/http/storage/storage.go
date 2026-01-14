package storage

type DB interface {
	CreateStudent(name string, email string, password string) (int, error)
}
