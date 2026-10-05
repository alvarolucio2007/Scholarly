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

// GetEnrollmentByID godoc
//
//	@Summary		Fetch an enrollment by it's ID
//	@Description	Fetch an enrollment by it's ID
//	@Tags			enrollments
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Enrollment ID"
//	@Success		200	{object}	EnrollmentResponseDTO
//	@Failure		400	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/enrollments/{id} [get]
func (h *Handler) GetEnrollmentByID(w http.ResponseWriter, r *http.Request) {
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

// ListEnrollments godoc
//
//	@Summary		List enrollments by parameters
//	@Description	List enrollments by courseID,status and studentID
//	@Tags			enrollments
//	@Accept			json
//	@Produce		json
//	@Param			course_id	query		int		false	"filter by course id"
//	@Param			status		query		string	false	"filter by status"
//	@Param			student_id	query		int		false	"filter by student id"
//	@Success		200			{array}		EnrollmentResponseDTO
//	@Failure		400			{object}	ErrorResponse
//	@Failure		500			{object}	ErrorResponse
//	@Router			/enrollments [get]
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
			_ = writeJSONError(w, http.StatusBadRequest, "student_id must be greater than 0")
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
			return
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

// UpdateEnrollment godoc
//
//	@Summary		Update an enrollment
//	@Description	Update an enrollment by its ID
//	@Tags			enrollments
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		UpdateEnrollmentDTO	true	"enrollment DTO"
//	@Success		201		{object}	EnrollmentResponseDTO
//	@Failure		400		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		409		{object}	ErrorResponse
//	@Failure		422		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/enrollments [put]
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

// DeleteEnrollment godoc
//
//	@Summary		Deletes an enrollment
//	@Description	Deletes an grade by ID
//	@Tags			enrollments
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Enrollment ID"
//	@Success		204
//	@Failure		400	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/enrollments/{id} [delete]
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
