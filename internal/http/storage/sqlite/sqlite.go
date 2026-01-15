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

func (s *Sqlite) GetStudents() ([]types.Student, error) {
	stmt, err := s.Db.Prepare(`SELECT id, name, email FROM students`)
	if err != nil {
		return []types.Student{}, nil
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []types.Student

	for rows.Next() {
		var student types.Student
		err := rows.Scan(&student.Id, &student.Name, &student.Email)
		if err != nil {
			return nil, err
		}
		students = append(students, student)
	}

	return students, nil

}

func (s *Sqlite) UpdateStudent(id int, feild string, newValue any) (types.Student, error) {

	query := fmt.Sprintf("UPDATE students set %s = ? WHERE id = ? RETURNING id,name,email,password", feild)

	stmt, err := s.Db.Prepare(query)
	if err != nil {
		return types.Student{}, err
	}

	defer stmt.Close()

	row := stmt.QueryRow(newValue, id)
	if err != nil {
		return types.Student{}, err
	}

	var updatedStudent types.Student
	err = row.Scan(&updatedStudent.Id, &updatedStudent.Name, &updatedStudent.Email, &updatedStudent.Password)
	if err != nil {
		return types.Student{}, err
	}
	return updatedStudent, nil
}

func (s *Sqlite) DeleteStudent(id int) (int, error) {
	query := fmt.Sprintf("DELETE FROM students WHERE id = ? RETURNING %d", id)
	stmt, err := s.Db.Prepare(query)
	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	row := stmt.QueryRow(id)
	if err != nil {
		return 0, err
	}

	var deletedId int
	err = row.Scan(&deletedId)
	if err != nil {
		return 0, err
	}

	return deletedId, nil

}
