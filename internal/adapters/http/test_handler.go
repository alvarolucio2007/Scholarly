package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/alvarolucio2007/Scholarly/internal/ports"
	"github.com/alvarolucio2007/Scholarly/internal/services"
	"github.com/go-chi/chi/v5"
)

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
	if v := q.Get("test_date"); v != "" {
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			_ = writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		filter.TestDate = &t
	}
	tests, err := h.tests.ListTests(r.Context(), filter)
	if err != nil {
		respondDomainError(w, err)
		return
	}
	_ = writeJSONData(w, http.StatusOK, &tests)
}

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
