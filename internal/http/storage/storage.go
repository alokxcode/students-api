package storage

import "github.com/alokxcode/students-api/internal/types"

type DB interface {
	CreateStudent(name string, email string, password string) (int, error)
	GetStudentById(id int) (types.Student, error)
	GetStudents() ([]types.Student, error)
	UpdateStudent(id int, student types.Student) (int, error)
	DeleteStudent(id int) (int, error)
}
