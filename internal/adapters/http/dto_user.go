package http

type CreateUserDTO struct {
	Name     string `json:"name" validate:"required,max=255"`
	CPF      string `json:"cpf" validate:"required,max=14"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
type UpdateUserDTO struct {
	ID       int64   `json:"id" validate:"required"`
	Name     *string `json:"name" validate:"max=255"`
	CPF      *string `json:"cpf" validate:"max=14"`
	Email    *string `json:"email" validate:"email"`
	Password *string `json:"password"`
}
