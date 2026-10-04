package http

type CreateTeacherDTO struct {
	UserID     int64  `json:"user_id" validate:"required,gte=1"`
	Department string `json:"department" validate:"required,max=100"`
}
type UpdateTeacherDTO struct {
	UserID     int64   `json:"user_id" validate:"required,gte=1"`
	Department *string `json:"department" validate:"max=100"`
}
type TeacherResponseDTO struct {
	UserID     int64  `json:"user_id"`
	Department string `json:"department"`
}
