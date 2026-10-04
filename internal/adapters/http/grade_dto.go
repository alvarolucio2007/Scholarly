package http

type CreateGradeDTO struct {
	TestID       int64   `json:"test_id" validate:"required,gte=1"`
	EnrollmentID int64   `json:"enrollment_id" validate:"required,gte=1"`
	Value        float64 `json:"value" validate:"required,gte=0,lte=10"`
}
type UpdateGradeDTO struct {
	ID           int64    `json:"id" validate:"required,gte=1"`
	TestID       *int64   `json:"test_id" validate:"gte=1"`
	EnrollmentID *int64   `json:"enrollment_id" validate:"gte=1"`
	Value        *float64 `json:"value" validate:"gte=0,lte=10"`
}
