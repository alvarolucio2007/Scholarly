package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/alvarolucio2007/Scholarly/internal/ports"
	"github.com/alvarolucio2007/Scholarly/internal/services"
	"github.com/go-chi/chi/v5"
)

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
