package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/alvarolucio2007/Scholarly/internal/ports"
	"github.com/alvarolucio2007/Scholarly/internal/services"
	"github.com/go-chi/chi/v5"
)

// CreateCourse godoc
//
//	@Summary		Create a course
//	@Description	Create a course
//	@Tags			courses
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateCourseDTO	true	"User DTO"
//	@Success		201		{object}	domain.Course
//	@Failure		400		{object}	error
//	@Failure		409		{object}	error
//	@Failure		422		{object}	error
//	@Failure		500		{object}	error
//	@Router			/courses [post]
func (h *Handler) CreateCourse(w http.ResponseWriter, r *http.Request) {
	var dto CreateCourseDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		_ = writeJSONError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if err := validate.Struct(dto); err != nil {
		_ = writeJSONError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	course, err := h.courses.CreateCourse(r.Context(), services.CreateCoursePayload{
		TeacherID:   dto.TeacherID,
		Name:        dto.Name,
		Code:        dto.Code,
		Semester:    dto.Semester,
		MaxStudents: dto.MaxStudents,
	})
	if err != nil {
		respondDomainError(w, err)
		return
	}
	_ = writeJSONData(w, http.StatusCreated, course)
}

// GetCourseByID godoc
//
//	@Summary		Fetch a course by its ID
//	@Description	Fetch a course by its ID
//	@Tags			courses
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Course ID"
//	@Success		200	{object}	domain.Course
//	@Failure		400	{object}	errorResponse
//	@Failure		404	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Router			/courses/{id} [get]
func (h *Handler) GetCourseByID(w http.ResponseWriter, r *http.Request) {
	courseID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		_ = writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	if courseID <= 0 {
		_ = writeJSONError(w, http.StatusBadRequest, "courseID must be at least 1")
		return
	}
	course, err := h.courses.GetByID(r.Context(), courseID)
	if err != nil {
		respondDomainError(w, err)
		return
	}
	_ = writeJSONData(w, http.StatusOK, course)
}

// ListCourses godoc
//
//	@Summary		List courses by parameters
//	@Description	List courses by name,cpf and email
//	@Tags			courses
//	@Accept			json
//	@Produce		json
//	@Param			teacher_id	query		int		false	"filter by teacher ID"
//	@Param			name		query		string	false	"filter by name"
//	@Param			code		query		string	false	"filter by code"
//	@Param			semester	query		string	false	"filter by semester"
//	@Success		200			{array}		domain.Course
//	@Failure		400			{object}	errorResponse
//	@Failure		500			{object}	errorResponse
//	@Router			/courses [get]
func (h *Handler) ListCourses(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filter := ports.CourseFilter{}
	if v := q.Get("teacher_id"); v != "" {
		vID, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			_ = writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		if vID <= 0 {
			_ = writeJSONError(w, http.StatusBadRequest, "teacherID must be greater than 0")
			return
		}
		filter.TeacherID = &vID
	}
	if v := q.Get("name"); v != "" {
		if len(v) > 255 {
			_ = writeJSONError(w, http.StatusBadRequest, "name should be less than 255 characters")
			return
		}
		filter.Name = &v
	}
	if v := q.Get("code"); v != "" {
		if len(v) > 50 {
			_ = writeJSONError(w, http.StatusBadRequest, "code should be less than 255 characters")
			return
		}
		filter.Code = &v
	}
	if v := q.Get("semester"); v != "" {
		if len(v) > 50 {
			_ = writeJSONError(w, http.StatusBadRequest, "code should be less than 255 characters")
			return
		}
		filter.Semester = &v
	}
	courses, err := h.courses.ListCourses(r.Context(), filter)
	if err != nil {
		respondDomainError(w, err)
		return
	}
	_ = writeJSONData(w, http.StatusOK, &courses)
}

// UpdateCourse godoc
//
//	@Summary		Update a course
//	@Description	Update a course
//	@Tags			courses
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		UpdateCourseDTO	true	"Course DTO"
//	@Success		201		{object}	domain.Course
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		409		{object}	error
//	@Failure		422		{object}	error
//	@Failure		500		{object}	error
//	@Router			/courses [put]
func (h *Handler) UpdateCourse(w http.ResponseWriter, r *http.Request) {
	var dto UpdateCourseDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		_ = writeJSONError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if err := validate.Struct(dto); err != nil {
		_ = writeJSONError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	course, err := h.courses.UpdateCourse(r.Context(), services.UpdateCoursePayload{
		ID:          dto.ID,
		TeacherID:   dto.TeacherID,
		Name:        dto.Name,
		Code:        dto.Code,
		Semester:    dto.Semester,
		MaxStudents: dto.MaxStudents,
	})
	if err != nil {
		respondDomainError(w, err)
		return
	}
	_ = writeJSONData(w, http.StatusOK, course)
}

// DeleteCourse godoc
//
//	@Summary		Deletes an course
//	@Description	Deletes an course by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Course ID"
//	@Success		204
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Router			/courses/{id} [delete]
func (h *Handler) DeleteCourse(w http.ResponseWriter, r *http.Request) {
	courseID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		_ = writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	if courseID <= 0 {
		_ = writeJSONError(w, http.StatusBadRequest, "courseID must be at least 1")
		return
	}
	if err := h.courses.DeleteCourse(r.Context(), courseID); err != nil {
		respondDomainError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// EnrollStudent godoc
//
//	@Summary		Enroll a student
//	@Description	Enroll a student
//	@Tags			courses
//	@Accept			json
//	@Produce		json
//	@Param			student_id	path		int	true	"Student ID"
//	@Success		201			{object}	int
//	@Failure		400			{object}	error
//	@Failure		409			{object}	error
//	@Failure		422			{object}	error
//	@Failure		500			{object}	error
//	@Router			/courses/enroll/{student_id} [post]
func (h *Handler) EnrollStudent(w http.ResponseWriter, r *http.Request) {
	studentID, err := strconv.ParseInt(chi.URLParam(r, "student_id"), 10, 64)
	if err != nil {
		_ = writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	if studentID <= 0 {
		_ = writeJSONError(w, http.StatusBadRequest, "studentID must be at least 1")
		return
	}
	courseID, err := strconv.ParseInt(chi.URLParam(r, "course_id"), 10, 64)
	if err != nil {
		_ = writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	if courseID <= 0 {
		_ = writeJSONError(w, http.StatusBadRequest, "courseID must be at least 1")
		return
	}
	enrollmentID, err := h.courses.EnrollStudent(r.Context(), studentID, courseID)
	if err != nil {
		respondDomainError(w, err)
		return
	}
	_ = writeJSONData(w, http.StatusOK, enrollmentID)
}
