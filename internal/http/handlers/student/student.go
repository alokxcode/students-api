package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

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
