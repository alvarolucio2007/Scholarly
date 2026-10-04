package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/alvarolucio2007/Scholarly/internal/ports"
	"github.com/alvarolucio2007/Scholarly/internal/services"
	"github.com/go-chi/chi/v5"
)

// CreateTest godoc
//
//	@Summary		Create a test
//	@Description	Create a test
//	@Tags			tests
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateTestDTO	true	"Test DTO"
//	@Success		201		{object}	domain.Test
//	@Failure		400		{object}	error
//	@Failure		409		{object}	error
//	@Failure		422		{object}	error
//	@Failure		500		{object}	error
//	@Router			/tests [post]
func (h *Handler) CreateTest(w http.ResponseWriter, r *http.Request) {
	var dto CreateTestDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		_ = writeJSONError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if err := validate.Struct(dto); err != nil {
		_ = writeJSONError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	test, err := h.tests.CreateTest(r.Context(), services.CreateTestPayload{
		CourseID: dto.CourseID,
		Name:     dto.Name,
		Weight:   dto.Weight,
		TestDate: dto.TestDate,
	})
	if err != nil {
		respondDomainError(w, err)
		return
	}
	_ = writeJSONData(w, http.StatusCreated, test)
}

// GetTestByID godoc
//
//	@Summary		Fetch a test by its ID
//	@Description	Fetch a test by its ID
//	@Tags			Test
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Test ID"
//	@Success		200	{object}	domain.Test
//	@Failure		400	{object}	errorResponse
//	@Failure		404	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Router			/tests/{id} [get]

func (h *Handler) GetTestByID(w http.ResponseWriter, r *http.Request) {
	testID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		_ = writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	if testID <= 0 {
		_ = writeJSONError(w, http.StatusBadRequest, "testID must be at least 1")
		return
	}
	test, err := h.tests.GetByID(r.Context(), testID)
	if err != nil {
		respondDomainError(w, err)
		return
	}
	_ = writeJSONData(w, http.StatusOK, test)
}

// ListTests godoc
//
//	@Summary		List tests by parameters
//	@Description	List tests by courseID and name
//	@Tags			Test
//	@Accept			json
//	@Produce		json
//	@Param			course_id	query		int		false	"filter by course ID"
//	@Param			name		query		string	false	"filter by name"
//	@Success		200			{array}		Test
//	@Failure		400			{object}	errorResponse
//	@Failure		500			{object}	errorResponse
//	@Router			/tests/{id} [get]
func (h *Handler) ListTests(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filter := ports.TestFilter{}
	if v := q.Get("course_id"); v != "" {
		vID, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			_ = writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		if vID <= 0 {
			_ = writeJSONError(w, http.StatusBadRequest, "course_id must be greater than 0")
			return
		}
		filter.CourseID = &vID
	}
	if v := q.Get("name"); v != "" {
		if len(v) > 255 {
			_ = writeJSONError(w, http.StatusBadRequest, "name should be less than 255 characters")
			return
		}
		filter.Name = &v
	}
	tests, err := h.tests.ListTests(r.Context(), filter)
	if err != nil {
		respondDomainError(w, err)
		return
	}
	_ = writeJSONData(w, http.StatusOK, &tests)
}

// UpdateTest godoc
//
//	@Summary		Update a test
//	@Description	Update a test
//	@Tags			tests
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		UpdateTestDTO	true	"Test DTO"
//	@Success		201		{object}	domain.Test
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		409		{object}	error
//	@Failure		422		{object}	error
//	@Failure		500		{object}	error
//	@Router			/tests [put]
func (h *Handler) UpdateTest(w http.ResponseWriter, r *http.Request) {
	var dto UpdateTestDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		_ = writeJSONError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if err := validate.Struct(dto); err != nil {
		_ = writeJSONError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	test, err := h.tests.UpdateTest(r.Context(), services.UpdateTestPayload{
		ID:       dto.ID,
		CourseID: dto.CourseID,
		Name:     dto.Name,
		Weight:   dto.Weight,
		TestDate: dto.TestDate,
	})
	if err != nil {
		respondDomainError(w, err)
		return
	}
	_ = writeJSONData(w, http.StatusOK, test)
}

// DeleteTest godoc
//
//	@Summary		Deletes a test
//	@Description	Deletes an test by ID
//	@Tags			Test
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Test ID"
//	@Success		204
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Router			/tests/{id} [delete]
func (h *Handler) DeleteTest(w http.ResponseWriter, r *http.Request) {
	testID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		_ = writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	if testID <= 0 {
		_ = writeJSONError(w, http.StatusBadRequest, "testID must be at least 1")
		return
	}
	if err := h.tests.DeleteTest(r.Context(), testID); err != nil {
		respondDomainError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
