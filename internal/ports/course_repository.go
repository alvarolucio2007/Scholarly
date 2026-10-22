package ports

import (
	"context"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
)

type CourseFilter struct {
	TeacherID *int64
	Name      *string
	Code      *string
	Semester  *string
}
type CourseRepository interface {
	Repository[domain.Course]
	List(ctx context.Context, filter CourseFilter) ([]domain.Course, error)
	EnrollStudent(ctx context.Context, studentID, courseID int64) (int64, error)
}
