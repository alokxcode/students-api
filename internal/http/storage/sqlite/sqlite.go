package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/alokxcode/students-api/internal/config"
	"github.com/alokxcode/students-api/internal/types"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	Db *sql.DB
}

func New(conf *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", conf.Strorage_Path)
	if err != nil {
		return nil, err
	}

	fmt.Println(conf.Strorage_Path)

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	email TEXT,
	password TEXT
	)`)

	if err != nil {
		return nil, err
	}

	return &Sqlite{
		Db: db,
	}, nil

}

func (s *Sqlite) CreateStudent(name string, email string, password string) (int, error) {
	stm, err := s.Db.Prepare(`INSERT INTO students(name, email, password) VALUES(?,?,?)`)
	if err != nil {
		return 0, err
	}
	defer stm.Close()

	result, err := stm.Exec(name, email, password)
	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(lastId), nil

}

func (s *Sqlite) GetStudentById(id int) (types.Student, error) {
	stm, err := s.Db.Prepare(`SELECT * FROM students WHERE id = ? LIMIT 1`)
	if err != nil {
		return types.Student{}, err
	}
	defer stm.Close()

	var student types.Student
	err = stm.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no student found")
		}
		return types.Student{}, err
	}

	return student, nil

}
