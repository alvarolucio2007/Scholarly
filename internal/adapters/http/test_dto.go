package http

import "time"

type CreateTestDTO struct {
	CourseID int64     `json:"course_id" validate:"required,gte=1"`
	Name     string    `json:"name" validate:"required, max=255"`
	Weight   float64   `json:"weight" validate:"required gt=0" `
	TestDate time.Time `json:"test_date" validate:"required,datetime"`
}
type UpdateTestDTO struct {
	ID       int64      `json:"id" validate:"required,gte=1"`
	CourseID *int64     `json:"course_id" validate:"gte=1"`
	Name     *string    `json:"name" validate:"required,gte=1"`
	Weight   *float64   `json:"weight" validate:"gt=0" `
	TestDate *time.Time `json:"test_date" validate:"datetime"`
}
