package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/alvarolucio2007/Scholarly/internal/ports"
	"github.com/alvarolucio2007/Scholarly/internal/services"
	"github.com/go-chi/chi/v5"
)

// CreateTeacher godoc
//
//	@Summary		Create an teacher account
//	@Description	Create an teacher account
//	@Tags			teachers
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateTeacherDTO	true	"User DTO"
//	@Success		201		{object}	TeacherResponseDTO
//	@Failure		400		{object}	error
//	@Failure		409		{object}	error
//	@Failure		422		{object}	error
//	@Failure		500		{object}	error
//	@Router			/teachers [post]

func (h *Handler) CreateTeacher(w http.ResponseWriter, r *http.Request) {
	var dto CreateTeacherDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		_ = writeJSONError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if err := validate.Struct(dto); err != nil {
		_ = writeJSONError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	teacher, err := h.teachers.CreateTeacher(r.Context(), services.CreateTeacherPayload{
		UserID:     dto.UserID,
		Department: dto.Department,
	})
	if err != nil {
		respondDomainError(w, err)
		return
	}
	_ = writeJSONData(w, http.StatusCreated, teacher)
}

// GetTeacherByID godoc
//
//	@Summary		Fetch a teacher by its ID
//	@Description	Fetch a teacher by its ID
//	@Tags			teachers
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Teacher ID"
//	@Success		200	{object}	TeacherResponseDTO
//	@Failure		400	{object}	errorResponse
//	@Failure		404	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Router			/teachers/{id} [get]
func (h *Handler) GetTeacherByID(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		_ = writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	if userID <= 0 {
		_ = writeJSONError(w, http.StatusBadRequest, "userID must be at least 1")
		return
	}
	user, err := h.teachers.GetByID(r.Context(), userID)
	if err != nil {
		respondDomainError(w, err)
		return
	}
	if err := writeJSONData(w, http.StatusOK, user); err != nil {
		_ = writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

// ListTeachers godoc
//
//	@Summary		List teachers by department
//	@Description	List teachers by department
//	@Tags			teachers
//	@Accept			json
//	@Produce		json
//	@Param			department	query		string	false	"filter by department"
//	@Success		200			{array}		UserResponseDTO
//	@Failure		400			{object}	errorResponse
//	@Failure		500			{object}	errorResponse
//	@Router			/teachers [get]
func (h *Handler) ListTeachers(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filter := ports.TeacherFilter{}
	if v := q.Get("department"); v != "" {
		if len(v) > 100 {
			_ = writeJSONError(w, http.StatusBadRequest, "department should be less than 100 characters long")
			return
		}
		filter.Department = &v
	}
	teachers, err := h.teachers.List(r.Context(), filter)
	if err != nil {
		respondDomainError(w, err)
		return
	}
	_ = writeJSONData(w, http.StatusOK, teachers)
}

// UpdateTeacher godoc
//
//	@Summary		Update an teacher account
//	@Description	Update an teacher account
//	@Tags			teachers
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		UpdateTeacherDTO	true	"Teacher DTO"
//	@Success		201		{object}	UserResponseDTO
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		409		{object}	error
//	@Failure		422		{object}	error
//	@Failure		500		{object}	error
//	@Router			/teachers [put]
func (h *Handler) UpdateTeacher(w http.ResponseWriter, r *http.Request) {
	var dto UpdateTeacherDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		_ = writeJSONError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if err := validate.Struct(dto); err != nil {
		_ = writeJSONError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	teacher, err := h.teachers.Update(r.Context(), services.UpdateTeacherPayload{
		UserID:     dto.UserID,
		Department: dto.Department,
	})
	if err != nil {
		respondDomainError(w, err)
		return
	}
	_ = writeJSONData(w, http.StatusOK, teacher)
}

// DeleteTeacher godoc
//
//	@Summary		Deletes an teacher account
//	@Description	Deletes an teacher account by ID
//	@Tags			teachers
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Teacher ID"
//	@Success		204
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Router			/teachers/{id} [delete]

func (h *Handler) DeleteTeacher(w http.ResponseWriter, r *http.Request) {
	teacherID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		_ = writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	if teacherID <= 0 {
		_ = writeJSONError(w, http.StatusBadRequest, "teacherID must be at least 1")
		return
	}
	if err := h.teachers.Delete(r.Context(), teacherID); err != nil {
		respondDomainError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
