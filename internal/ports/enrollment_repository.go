package ports

import (
	"context"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
)

type EnrollmentFilter struct {
	StudentID *int64
	CourseID  *int64
	Status    *string
}
type EnrollmentRepository interface {
	Repository[domain.Enrollment] // TODO: Add status editing later for Update
	List(ctx context.Context, filter EnrollmentFilter) ([]*domain.Enrollment, error)
	GetWeightedAverage(ctx context.Context, studentID, courseID int64) (float64, error)
}
