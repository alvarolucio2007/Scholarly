package http

import "github.com/alvarolucio2007/Scholarly/internal/domain"

type CreateEnrollmentDTO struct {
	StudentID int64                   `json:"student_id" validate:"required,gte=1"`
	CourseID  int64                   `json:"course_id" validate:"required,gte=1"`
	Status    domain.EnrollmentStatus `json:"status" validate:"required"`
}
type UpdateEnrollmentDTO struct {
	ID        int64  `json:"id" validate:"gte=1"`
	StudentID *int64 `json:"student_id" validate:"gte=1"`
	CourseID  *int64 `json:"course_id" validate:"gte=1"`
}
