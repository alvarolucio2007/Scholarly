package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/alvarolucio2007/Scholarly/internal/ports"
	"github.com/alvarolucio2007/Scholarly/internal/services"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) GetErollmentByID(w http.ResponseWriter, r *http.Request) {
	courseID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		_ = writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	if courseID <= 0 {
		_ = writeJSONError(w, http.StatusBadRequest, "courseID must be at least 1")
		return
	}
	enrollment, err := h.enrollments.GetByID(r.Context(), courseID)
	if err != nil {
		respondDomainError(w, err)
		return
	}
	_ = writeJSONData(w, http.StatusOK, enrollment)
}

func (h *Handler) ListEnrollments(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filter := ports.EnrollmentFilter{}
	if v := q.Get("student_id"); v != "" {
		vID, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			_ = writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		if vID <= 0 {
			_ = writeJSONError(w, http.StatusBadRequest, "teacherID must be greater than 0")
			return
		}
		filter.StudentID = &vID
	}
	if v := q.Get("course_id"); v != "" {
		vID, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			_ = writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		if vID <= 0 {
			_ = writeJSONError(w, http.StatusBadRequest, "teacherID must be greater than 0")
			return
		}
		filter.CourseID = &vID
	}
	if v := q.Get("status"); v != "" {
		if v != string(domain.StatusActive) && v != string(domain.StatusCancelled) && v != string(domain.StatusCompleted) {
			_ = writeJSONError(w, http.StatusBadRequest, "status must be either 'active', 'cancelled' or 'completed' ")
		}
		filter.Status = &v
	}
	enrollments, err := h.enrollments.ListEnrollments(r.Context(), filter)
	if err != nil {
		respondDomainError(w, err)
		return
	}
	_ = writeJSONData(w, http.StatusOK, &enrollments)
}

func (h *Handler) UpdateEnrollment(w http.ResponseWriter, r *http.Request) {
	var dto UpdateEnrollmentDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		_ = writeJSONError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if err := validate.Struct(dto); err != nil {
		_ = writeJSONError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	enrollment, err := h.enrollments.UpdateEnrollment(r.Context(), services.UpdateEnrollmentPayload{
		ID:        dto.ID,
		StudentID: dto.StudentID,
		CourseID:  dto.CourseID,
	})
	if err != nil {
		respondDomainError(w, err)
		return
	}
	_ = writeJSONData(w, http.StatusOK, enrollment)
}

func (h *Handler) DeleteEnrollment(w http.ResponseWriter, r *http.Request) {
	enrollmentID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		_ = writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	if enrollmentID <= 0 {
		_ = writeJSONError(w, http.StatusBadRequest, "enrollmentID must be at least 1")
		return
	}
	if err := h.enrollments.DeleteEnrollment(r.Context(), enrollmentID); err != nil {
		respondDomainError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
