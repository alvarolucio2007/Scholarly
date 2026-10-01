package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	buf, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(buf); err != nil {
		return err
	}
	return nil
}

func writeJSONError(w http.ResponseWriter, status int, message string) error {
	type envelope struct {
		Error string `json:"error"`
	}
	return writeJSON(w, status, &envelope{Error: message})
}

func writeJSONData(w http.ResponseWriter, status int, data any) error {
	type envelope struct {
		Data any `json:"data"`
	}
	return writeJSON(w, status, &envelope{Data: data})
}

func respondDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		_ = writeJSONError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrConflict):
		_ = writeJSONError(w, http.StatusConflict, err.Error())
	case errors.Is(err, domain.ErrEmailAlreadyExists):
		_ = writeJSONError(w, http.StatusConflict, err.Error())
	case errors.Is(err, domain.ErrCPFAlreadyExists):
		_ = writeJSONError(w, http.StatusConflict, err.Error())
	case errors.Is(err, domain.ErrStudentNotFound):
		_ = writeJSONError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrCourseNotFound):
		_ = writeJSONError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrEnrollmentAlreadyExists):
		_ = writeJSONError(w, http.StatusConflict, err.Error())
	case errors.Is(err, domain.ErrCourseFull):
		_ = writeJSONError(w, http.StatusConflict, err.Error())
	default:
		_ = writeJSONError(w, http.StatusInternalServerError, "internal error")
	}
}
