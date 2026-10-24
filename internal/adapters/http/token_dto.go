package http

type CreateUserTokenPayloadDTO struct {
	Email    string `json:"email" validate:"required,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}
