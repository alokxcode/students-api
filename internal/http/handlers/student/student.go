package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/alokxcode/students-api/internal/http/storage"
	"github.com/alokxcode/students-api/internal/types"
	"github.com/alokxcode/students-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New(db storage.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Creating student")

		var student types.Student

		// decode
		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GenerelError(err, http.StatusBadRequest))
			return
		}

		// validate
		if errs := validator.New().Struct(student); errs != nil {
			validateErr := errs.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErr, http.StatusBadRequest))
			return
		}

		// create
		lastId, err := db.CreateStudent(
			student.Name,
			student.Email,
			student.Password,
		)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
		}

		slog.Info("User created successfully", slog.String("Id :", fmt.Sprint(lastId)))

		response.WriteJson(w, http.StatusCreated, map[string]int{"id": lastId})
	}
}

func GetById(db storage.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("getting student", slog.String("id :", id))
		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, err)
		}

		student, err := db.GetStudentById(int(intId))
		if err != nil {
			response.GenerelError(err, 400)
			return
		}

		response.WriteJson(w, http.StatusOK, student)

	}
}

func GetList(db storage.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Gettings students")
		students, err := db.GetStudents()
		if err != nil {
			response.GenerelError(err, 500)
			return
		}

		response.WriteJson(w, http.StatusOK, students)
	}
}

func Update(db storage.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		IntId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.GenerelError(err, 400)
			return
		}

		slog.Info("updating student", slog.String("id : ", fmt.Sprint(IntId)))

		var student types.StudentPatch
		// decode
		err = json.NewDecoder(r.Body).Decode(&student)
		if err != nil {
			response.GenerelError(err, 400)
			return
		}

		var updatedStudent types.Student

		if student.Name != nil {
			updatedStudent, err = db.UpdateStudent(int(IntId), "name", student.Name)
			if err != nil {
				response.GenerelError(err, 500)
				return
			}

		}

		if student.Email != nil {
			updatedStudent, err = db.UpdateStudent(int(IntId), "email", student.Email)
			if err != nil {
				response.GenerelError(err, 500)
				return
			}

		}

		if student.Password != nil {
			updatedStudent, err = db.UpdateStudent(int(IntId), "password", student.Password)
			if err != nil {
				response.GenerelError(err, 500)
				return
			}
		}
		slog.Info("student updated", slog.String("id : ", fmt.Sprint(IntId)))
		response.WriteJson(w, http.StatusOK, updatedStudent)

	}
}

func Delete(db storage.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		IntId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.GenerelError(err, 400)
			return
		}

		slog.Info("deleting student", slog.String("", fmt.Sprintf("Id : %d", IntId)))

		deltedId, err := db.DeleteStudent(int(IntId))
		if err != nil && deltedId == 0 {
			response.GenerelError(err, 500)
			return
		}

		response.WriteJson(w, http.StatusOK, fmt.Sprintf("Deleted Student Id : %d", deltedId))

	}
}
