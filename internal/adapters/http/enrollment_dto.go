package http

import "time"

type CreateEnrollmentDTO struct {
	StudentID int64 `json:"student_id" validate:"required,gte=1"`
	CourseID  int64 `json:"course_id" validate:"required,gte=1"`
}
type UpdateEnrollmentDTO struct {
	ID        int64  `json:"id" validate:"gte=1"`
	StudentID *int64 `json:"student_id" validate:"gte=1"`
	CourseID  *int64 `json:"course_id" validate:"gte=1"`
}
type EnrollmentResponseDTO struct {
	ID         int64     `json:"id"`
	StudentID  int64     `json:"student_id"`
	CourseID   int64     `json:"course_id"`
	EnrolledAt time.Time `json:"enrolled_at"`
	Status     string    `json:"status"`
}
