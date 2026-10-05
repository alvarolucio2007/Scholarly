package http

import "time"

type CreateTestDTO struct {
	CourseID int64   `json:"course_id" validate:"required,gte=1"`
	Name     string  `json:"name" validate:"required,max=255"`
	Weight   float64 `json:"weight" validate:"required,gt=0" `
	TestDate string  `json:"test_date" validate:"required" example:"2026-03-15T00:00:00Z"`
}
type UpdateTestDTO struct {
	ID       int64    `json:"id" validate:"required,gte=1"`
	CourseID *int64   `json:"course_id" validate:"gte=1"`
	Name     *string  `json:"name" validate:"gte=1"`
	Weight   *float64 `json:"weight" validate:"gte=0" `
	TestDate *string  `json:"test_date" example:"2026-03-15T00:00:00Z"`
}
type TestResponseDTO struct {
	ID        int64      `json:"id"`
	CourseID  int64      `json:"course_id"`
	Name      string     `json:"name"`
	Weight    float64    `json:"weight"`
	TestDate  time.Time  `json:"test_date"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
