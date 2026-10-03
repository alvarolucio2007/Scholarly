package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/alvarolucio2007/Scholarly/internal/ports"
	"github.com/alvarolucio2007/Scholarly/internal/services"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) CreateGrade(w http.ResponseWriter, r *http.Request) {
	var dto CreateGradeDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		_ = writeJSONError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if err := validate.Struct(dto); err != nil {
		_ = writeJSONError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	grade, err := h.grades.CreateGrade(r.Context(), services.CreateGradePayload{
		TestID:       dto.TestID,
		EnrollmentID: dto.EnrollmentID,
		Value:        dto.Value,
	})
	if err != nil {
		respondDomainError(w, err)
		return
	}
	_ = writeJSONData(w, http.StatusCreated, grade)
}

func (h *Handler) GetGradeByID(w http.ResponseWriter, r *http.Request) {
	gradeID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		_ = writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	if gradeID <= 0 {
		_ = writeJSONError(w, http.StatusBadRequest, "gradeID must be at least 1")
		return
	}
	grade, err := h.grades.GetByID(r.Context(), gradeID)
	if err != nil {
		respondDomainError(w, err)
		return
	}
	_ = writeJSONData(w, http.StatusOK, grade)
}

func (h *Handler) ListGrades(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filter := ports.GradeFilter{}
	if v := q.Get("test_id"); v != "" {
		vID, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			_ = writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		if vID <= 0 {
			_ = writeJSONError(w, http.StatusBadRequest, "test_id must be greater than 0")
			return
		}
		filter.TestID = &vID
	}
	if v := q.Get("enrollment_id"); v != "" {
		vID, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			_ = writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		if vID <= 0 {
			_ = writeJSONError(w, http.StatusBadRequest, "enrollment_id must be greater than 0")
			return
		}
		filter.EnrollmentID = &vID
	}
	grades, err := h.grades.List(r.Context(), filter)
	if err != nil {
		respondDomainError(w, err)
		return
	}
	_ = writeJSONData(w, http.StatusOK, &grades)
}

func (h *Handler) UpdateGrade(w http.ResponseWriter, r *http.Request) {
	var dto UpdateGradeDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		_ = writeJSONError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if err := validate.Struct(dto); err != nil {
		_ = writeJSONError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	grade, err := h.grades.Update(r.Context(), services.UpdateGradePayload{
		ID:           dto.ID,
		TestID:       dto.TestID,
		EnrollmentID: dto.EnrollmentID,
		Value:        dto.Value,
	})
	if err != nil {
		respondDomainError(w, err)
		return
	}
	_ = writeJSONData(w, http.StatusOK, grade)
}

func (h *Handler) DeleteGrade(w http.ResponseWriter, r *http.Request) {
	gradeID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		_ = writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	if gradeID <= 0 {
		_ = writeJSONError(w, http.StatusBadRequest, "gradeID must be at least 1")
		return
	}
	if err := h.grades.Delete(r.Context(), gradeID); err != nil {
		respondDomainError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListReportCard(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filter := ports.ReportCardFilter{}
	if v := q.Get("student_id"); v != "" {
		vID, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			_ = writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		if vID <= 0 {
			_ = writeJSONError(w, http.StatusBadRequest, "student_id must be greater than 0")
			return
		}
		filter.StudentID = &vID
	}
	reportCards, err := h.grades.ListReportCard(r.Context(), filter)
	if err != nil {
		respondDomainError(w, err)
		return
	}
	_ = writeJSONData(w, http.StatusOK, &reportCards)
}
