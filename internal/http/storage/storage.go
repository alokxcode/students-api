package storage

import "github.com/alokxcode/students-api/internal/types"

type DB interface {
	CreateStudent(name string, email string, password string) (int, error)
	GetStudentById(id int) (types.Student, error)
	GetStudents() ([]types.Student, error)
	UpdateStudent(id int, feild string, newValue any) (types.Student, error)
	DeleteStudent(id int) (int, error)
}
