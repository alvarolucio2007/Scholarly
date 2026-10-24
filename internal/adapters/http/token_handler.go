package http

import (
	"encoding/json"
	"net/http"
)

// createTokenHandler godoc
//
//	@Summary		Creates a token
//	@Description	Creates a token for the user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateUserTokenPayloadDTO	true	"User credentials"
//	@Success		200		{object}	string						"token registered"
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/token [post]
func (h *Handler) createTokenHandler(w http.ResponseWriter, r *http.Request) {
	var dto CreateUserTokenPayloadDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		_ = writeJSONError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if err := validate.Struct(dto); err != nil {
		_ = writeJSONError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	token, err := h.users.VerifyUser(r.Context(), dto.Email, dto.Password)
	if err != nil {
		respondDomainError(w, err)
		return
	}
	_ = writeJSONData(w, http.StatusOK, map[string]string{"token": token})
}
