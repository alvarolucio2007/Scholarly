package http

type CreateStudentDTO struct {
	UserID           int64  `json:"user_id" validate:"required,gte=1"`
	EnrollmentNumber string `json:"enrollment_number" validate:"required" `
}
type UpdateStudentDTO struct {
	UserID           int64   `json:"user_id" validate:"required,gte=1"`
	EnrollmentNumber *string `json:"enrollment_number" `
}
