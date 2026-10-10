package ports

import (
	"context"
	"time"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
)

type Repository[T any] interface {
	Create(ctx context.Context, item *T) error
	GetByID(ctx context.Context, id int64) (*T, error)
	Update(ctx context.Context, item *T) error
	Delete(ctx context.Context, id int64) error
}
type UserFilter struct {
	Name  *string
	CPF   *string
	Email *string
}
type UserRepository interface {
	Repository[domain.User]
	List(ctx context.Context, filter UserFilter) ([]*domain.User, error)
}

type StudentRepository interface {
	Repository[domain.Student]
	GetByEnrollment(ctx context.Context, enrollment string) (*domain.Student, error)
}

type TeacherFilter struct {
	Department *string
}
type TeacherRepository interface {
	Repository[domain.Teacher]
	List(ctx context.Context, filter TeacherFilter) ([]*domain.Teacher, error)
}

type CourseFilter struct {
	TeacherID *int64
	Name      *string
	Code      *string
	Semester  *string
}
type CourseRepository interface {
	Repository[domain.Course]
	List(ctx context.Context, filter CourseFilter) ([]*domain.Course, error)
	EnrollStudent(ctx context.Context, studentID, courseID int64) (int64, error)
}

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

type TestFilter struct {
	CourseID *int64
	Name     *string
	TestDate *time.Time
}
type TestRepository interface {
	Repository[domain.Test]
	List(ctx context.Context, filter TestFilter) ([]*domain.Test, error)
}

type GradeFilter struct {
	TestID       *int64
	EnrollmentID *int64
}
type ReportCardFilter struct {
	StudentID *int64
}
type GradeRepository interface {
	Repository[domain.Grade]
	List(ctx context.Context, filter GradeFilter) ([]*domain.Grade, error)
	ListAverages(ctx context.Context, filter ReportCardFilter) ([]*domain.ReportCard, error)
	GetAverage(ctx context.Context, studentID, courseID int64) (*float64, error)
}
