package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

func WriteJson(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)

}

type ErrorResponse struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}

func GenerelError(err error, status int) ErrorResponse {
	return ErrorResponse{
		Status: status,
		Error:  err.Error(),
	}

}

func ValidationError(errs validator.ValidationErrors, status int) ErrorResponse {
	var errMsg []string
	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsg = append(errMsg, fmt.Sprintf("%s feild is requried", err.Field()))
		default:
			errMsg = append(errMsg, "invalid feild")
		}
	}
	return ErrorResponse{
		Status: status,
		Error:  strings.Join(errMsg, ","),
	}
}
