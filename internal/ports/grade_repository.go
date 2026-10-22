package ports

import (
	"context"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
)

type GradeFilter struct {
	TestID       *int64
	EnrollmentID *int64
}
type ReportCardFilter struct {
	StudentID *int64
}
type GradeRepository interface {
	Repository[domain.Grade]
	List(ctx context.Context, filter GradeFilter) ([]domain.Grade, error)
	ListAverages(ctx context.Context, filter ReportCardFilter) ([]domain.ReportCard, error)
	GetAverage(ctx context.Context, studentID, courseID int64) (float64, error)
}
