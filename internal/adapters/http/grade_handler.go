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

func newGradeResponseDTO(grade domain.Grade) GradeResponseDTO {
	return GradeResponseDTO{ID: grade.ID, TestID: grade.TestID, EnrollmentID: grade.EnrollmentID, Value: grade.Value, CreatedAt: grade.CreatedAt, UpdatedAt: grade.UpdatedAt}
}

// CreateGrade godoc
//
//	@Summary		Create a grade
//	@Description	Create a grade
//	@Tags			grades
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateGradeDTO	true	"Grade DTO"
//	@Success		201		{object}	GradeResponseDTO
//	@Failure		400		{object}	ErrorResponse
//	@Failure		409		{object}	ErrorResponse
//	@Failure		422		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/grades [post]
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
	grade, err := h.grades.Create(r.Context(), services.CreateGradePayload{
		TestID:       dto.TestID,
		EnrollmentID: dto.EnrollmentID,
		Value:        dto.Value,
	})
	if err != nil {
		respondDomainError(w, err)
		return
	}
	response := newGradeResponseDTO(grade)
	_ = writeJSONData(w, http.StatusCreated, response)
}

// GetGradeByID godoc
//
//	@Summary		Fetch a grade by its ID
//	@Description	Fetch a grade by its ID
//	@Tags			grades
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Grade ID"
//	@Success		200	{object}	GradeResponseDTO
//	@Failure		400	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/grades/{id} [get]
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
	response := newGradeResponseDTO(grade)
	_ = writeJSONData(w, http.StatusOK, response)
}

func newListGradeResponseDTO(grades []domain.Grade) ListGradeResponseDTO {
	responses := make([]GradeResponseDTO, 0, len(grades))
	for _, grade := range grades {
		responses = append(responses, newGradeResponseDTO(grade))
	}
	return ListGradeResponseDTO{
		Data: responses,
	}
}

// ListGrades godoc
//
//	@Summary		List grades by parameters
//	@Description	List grades by test_id and enrollment_id
//	@Tags			grades
//	@Accept			json
//	@Produce		json
//	@Param			test_id			query		int	false	"filter by test ID"
//	@Param			enrollment_id	query		int	false	"filter by enrollment ID"
//	@Success		200				{array}		GradeResponseDTO
//	@Failure		400				{object}	ErrorResponse
//	@Failure		500				{object}	ErrorResponse
//	@Router			/grades [get]
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
	response := newListGradeResponseDTO(grades)
	_ = writeJSONData(w, http.StatusOK, response)
}

// UpdateGrade godoc
//
//	@Summary		Update a grade
//	@Description	Update a grade
//	@Tags			grades
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		UpdateGradeDTO	true	"Grade DTO"
//	@Success		204
//	@Failure		400		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		409		{object}	ErrorResponse
//	@Failure		422		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/grades [put]
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
	err := h.grades.Update(r.Context(), services.UpdateGradePayload{
		ID:           dto.ID,
		TestID:       nilIfZeroInt(dto.TestID),
		EnrollmentID: nilIfZeroInt(dto.EnrollmentID),
		Value:        nilIfZeroFloat(dto.Value),
	})
	if err != nil {
		respondDomainError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DeleteGrade godoc
//
//	@Summary		Deletes a grade
//	@Description	Deletes a grade by ID
//	@Tags			grades
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Grade ID"
//	@Success		204
//	@Failure		400	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/grades/{id} [delete]
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

//TODO: Implement DTO for both handlers below

// ListAverages godoc
//
//	@Summary		List averages
//	@Description	List averages by studentID
//	@Tags			grades
//	@Accept			json
//	@Produce		json
//	@Param			student_id	query		int	false	"filter by studentID"
//	@Success		200			{array}		ListAverageResponseDTO
//	@Failure		400			{object}	ErrorResponse
//	@Failure		500			{object}	ErrorResponse
//	@Router			/grades/average/list [get]
func (h *Handler) ListAverages(w http.ResponseWriter, r *http.Request) {
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
	reportCards, err := h.grades.ListAverages(r.Context(), filter)
	if err != nil {
		respondDomainError(w, err)
		return
	}
	_ = writeJSONData(w, http.StatusOK, reportCards)
}

// GetAverage godoc
//
//	@Summary		Gets an student average
//	@Description	Gets an student average by userID and courseID
//	@Tags			grades
//	@Accept			json
//	@Produce		json
//	@Param			student_id	path		int	true	"Student ID"
//	@Param			course_id	path		int	true	"Course ID"
//	@Success		200			{object}	float64
//	@Failure		400			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Failure		500			{object}	ErrorResponse
//	@Router			/grades/average/{student_id}/{course_id} [get]
func (h *Handler) GetAverage(w http.ResponseWriter, r *http.Request) {
	studentID, err := strconv.ParseInt(chi.URLParam(r, "student_id"), 10, 64)
	if err != nil {
		_ = writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	if studentID <= 0 {
		_ = writeJSONError(w, http.StatusBadRequest, "student_id must be at least 1")
		return
	}
	courseID, err := strconv.ParseInt(chi.URLParam(r, "course_id"), 10, 64)
	if err != nil {
		_ = writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	if courseID <= 0 {
		_ = writeJSONError(w, http.StatusBadRequest, "course_id must be at least 1")
		return
	}
	res, err := h.grades.GetAverage(r.Context(), studentID, courseID)
	if err != nil {
		respondDomainError(w, err)
		return
	}
	_ = writeJSONData(w, 200, res)
}
