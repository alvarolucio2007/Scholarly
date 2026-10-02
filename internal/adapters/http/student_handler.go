package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/alvarolucio2007/Scholarly/internal/services"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) CreateStudent(w http.ResponseWriter, r *http.Request) {
	var dto CreateStudentDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		_ = writeJSONError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if err := validate.Struct(dto); err != nil {
		_ = writeJSONError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	student, err := h.students.CreateStudent(r.Context(), services.CreateStudentPayload{
		UserID:           dto.UserID,
		EnrollmentNumber: dto.EnrollmentNumber,
	})
	if err != nil {
		respondDomainError(w, err)
		return
	}
	_ = writeJSONData(w, http.StatusCreated, student)
}

func (h *Handler) GetStudentByID(w http.ResponseWriter, r *http.Request) {
	studentID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		_ = writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	if studentID <= 0 {
		_ = writeJSONError(w, http.StatusBadRequest, "userID must be at least 1")
		return
	}
	student, err := h.students.GetByID(r.Context(), studentID)
	if err != nil {
		respondDomainError(w, err)
		return
	}
	if err := writeJSONData(w, http.StatusOK, student); err != nil {
		_ = writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *Handler) GetStudentByEnrollment(w http.ResponseWriter, r *http.Request) {
	studentEnrollment := chi.URLParam(r, "enrollment")
	if studentEnrollment == "" {
		_ = writeJSONError(w, http.StatusBadRequest, "studentEnrollment is missing")
		return
	}
	student, err := h.students.GetByEnrollment(r.Context(), studentEnrollment)
	if err != nil {
		respondDomainError(w, err)
		return
	}
	if err := writeJSONData(w, http.StatusOK, student); err != nil {
		_ = writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *Handler) UpdateStudent(w http.ResponseWriter, r *http.Request) {
	var dto UpdateStudentDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		_ = writeJSONError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if err := validate.Struct(dto); err != nil {
		_ = writeJSONError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	student, err := h.students.Update(r.Context(), services.UpdateStudentPayload{
		UserID:           dto.UserID,
		EnrollmentNumber: dto.EnrollmentNumber,
	})
	if err != nil {
		respondDomainError(w, err)
		return
	}
	_ = writeJSONData(w, http.StatusOK, student)
}

func (h *Handler) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	studentID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		_ = writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	if studentID <= 0 {
		_ = writeJSONError(w, http.StatusBadRequest, "studentID must be at least 1")
		return
	}
	if err := h.students.Delete(r.Context(), studentID); err != nil {
		respondDomainError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
