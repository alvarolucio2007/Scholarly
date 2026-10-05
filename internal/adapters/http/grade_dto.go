package http

import "time"

type CreateGradeDTO struct {
	TestID       int64   `json:"test_id" validate:"required,gte=1"`
	EnrollmentID int64   `json:"enrollment_id" validate:"required,gte=1"`
	Value        float64 `json:"value" validate:"required,gte=0,lte=10"`
}
type UpdateGradeDTO struct {
	ID           int64   `json:"id" validate:"gte=1"`
	TestID       int64   `json:"test_id" validate:"gte=0"`
	EnrollmentID int64   `json:"enrollment_id" validate:"gte=0"`
	Value        float64 `json:"value" validate:"lte=10" example:"0"`
}
type GradeResponseDTO struct {
	ID           int64      `json:"id"`
	TestID       int64      `json:"test_id"`
	EnrollmentID int64      `json:"enrollment_id"`
	Value        float64    `json:"value"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}
type ReportCardResponseDTO struct {
	StudentID   int64   `json:"student_id"`
	CourseID    int64   `json:"course_id"`
	StudentName string  `json:"student_name"`
	CourseName  string  `json:"course_name"`
	CourseCode  string  `json:"course_code"`
	TeacherName string  `json:"teacher_name"`
	Average     float64 `json:"average"`
	Situation   string  `json:"situation"`
}
