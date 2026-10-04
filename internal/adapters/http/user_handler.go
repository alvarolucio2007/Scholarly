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

func newUserResponse(user *domain.User) *UserResponseDTO {
	return &UserResponseDTO{ID: user.ID, Name: user.Name, CPF: user.CPF, Email: user.Email, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt}
}

// CreateUser godoc
//
//	@Summary		Create an user account
//	@Description	Create an user account
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateUserDTO	true	"User DTO"
//	@Success		201		{object}	UserResponseDTO
//	@Failure		400		{object}	error
//	@Failure		409		{object}	error
//	@Failure		422		{object}	error
//	@Failure		500		{object}	error
//	@Router			/users [post]
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
	response := newUserResponse(user)
	_ = writeJSONData(w, http.StatusCreated, response)
}

// GetUserByID godoc
//
//	@Summary		Fetch an user by it's ID
//	@Description	Fetch an user by it's ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	UserResponseDTO
//	@Failure		400	{object}	errorResponse
//	@Failure		404	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Router			/users/{id} [get]
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
	response := newUserResponse(user)
	_ = writeJSONData(w, http.StatusOK, response)
}

// ListUsers godoc
//
//	@Summary		List users by parameters
//	@Description	List users by name,cpf and email
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			name	query		string	false	"filter by name"
//	@Param			cpf		query		string	false	"filter by cpf"
//	@Param			email	query		string	false	"filter by email"
//	@Success		200		{array}		UserResponseDTO
//	@Failure		400		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Router			/users/ [get]
func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
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
	responses := make([]*UserResponseDTO, 0, len(users))
	for _, u := range users {
		response := newUserResponse(u)
		responses = append(responses, response)
	}
	_ = writeJSONData(w, http.StatusOK, responses)
}

// UpdateUser godoc
//
//	@Summary		Update an user account
//	@Description	Update an user account
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		UpdateUserDTO	true	"User DTO"
//	@Success		201		{object}	UserResponseDTO
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		409		{object}	error
//	@Failure		422		{object}	error
//	@Failure		500		{object}	error
//	@Router			/users [put]
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
	updateUserPayload := services.UpdateUserPayload{
		ID:       dto.ID,
		Name:     nilIfEmpty(dto.Name),
		CPF:      nilIfEmpty(dto.CPF),
		Email:    nilIfEmpty(dto.Email),
		Password: nilIfEmpty(dto.Password),
	}
	user, err := h.users.UpdateUser(r.Context(), updateUserPayload)
	if err != nil {
		respondDomainError(w, err)
		return
	}
	response := newUserResponse(user)
	_ = writeJSONData(w, http.StatusOK, response)
}

// DeleteUser godoc
//
//	@Summary		Deletes an user account
//	@Description	Deletes an user account by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"User ID"
//	@Success		204
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Router			/users/{id} [delete]
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
