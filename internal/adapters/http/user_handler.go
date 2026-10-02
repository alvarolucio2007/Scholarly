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

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var dto CreateUserDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		_ = writeJSONError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if err := validate.Struct(dto); err != nil {
		_ = writeJSONError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	user, err := h.users.CreateUser(r.Context(), services.CreateUserPayload{
		Name:     dto.Name,
		CPF:      dto.CPF,
		Email:    dto.Email,
		Password: dto.Password,
	})
	if err != nil {
		respondDomainError(w, err)
		return
	}
	_ = writeJSONData(w, http.StatusCreated, user)
}

func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		_ = writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	if userID <= 0 {
		_ = writeJSONError(w, http.StatusBadRequest, "userID must be at least 1")
		return
	}
	user, err := h.users.GetUserByID(r.Context(), userID)
	if err != nil {
		respondDomainError(w, err)
		return
	}
	_ = writeJSONData(w, http.StatusOK, user)
}

func (h *Handler) ListUser(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filter := ports.UserFilter{}
	if v := q.Get("name"); v != "" {
		if len(v) > 255 {
			_ = writeJSONError(w, http.StatusBadRequest, "name should be less than 255 characters long")
			return
		}
		filter.Name = &v
	}
	if v := q.Get("cpf"); v != "" {
		user := domain.User{CPF: v}
		if err := user.ValidateCPF(); err != nil {
			_ = writeJSONError(w, http.StatusBadRequest, "CPF is invalid")
			return
		}
		filter.CPF = &v
	}
	if v := q.Get("email"); v != "" {
		filter.Email = &v
	}
	users, err := h.users.ListUsers(r.Context(), filter)
	if err != nil {
		respondDomainError(w, err)
		return
	}
	_ = writeJSONData(w, http.StatusOK, &users)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var dto UpdateUserDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		_ = writeJSONError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if err := validate.Struct(dto); err != nil {
		_ = writeJSONError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	user, err := h.users.UpdateUser(r.Context(), services.UpdateUserPayload{
		ID:       dto.ID,
		Name:     dto.Name,
		CPF:      dto.CPF,
		Email:    dto.Email,
		Password: dto.Password,
	})
	if err != nil {
		respondDomainError(w, err)
		return
	}
	_ = writeJSONData(w, http.StatusOK, user)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		_ = writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	if userID <= 0 {
		_ = writeJSONError(w, http.StatusBadRequest, "userID must be at least 1")
		return
	}
	if err := h.users.Delete(r.Context(), userID); err != nil {
		respondDomainError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
