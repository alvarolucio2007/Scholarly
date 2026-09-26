package ports

import (
	"context"
	"time"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
)

type UserFilter struct {
	Name  *string
	CPF   *string
	Email *string
}
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, userID int64) (*domain.User, error)
	List(ctx context.Context, filter UserFilter) ([]*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, userID int64) error
}

type StudentRepository interface {
	Create(ctx context.Context, student *domain.Student) error
	GetByID(ctx context.Context, studentID int64) (*domain.Student, error)
	GetByEnrollment(ctx context.Context, enrollment string) (*domain.Student, error)
	Update(ctx context.Context, student *domain.Student) error
	Delete(ctx context.Context, studentID int64) error
}

type TeacherFilter struct {
	Department *string
}
type TeacherRepository interface {
	Create(ctx context.Context, teacher *domain.Teacher) error
	GetByID(ctx context.Context, teacherID int64) (*domain.Teacher, error)
	List(ctx context.Context, filter TeacherFilter) ([]*domain.Teacher, error)
	Update(ctx context.Context, teacher *domain.Teacher) error
	Delete(ctx context.Context, teacherID int64) error
}

type CourseFilter struct {
	TeacherID *int64
	Name      *string
	Code      *string
	Semester  *string
}
type CourseRepository interface {
	Create(ctx context.Context, course *domain.Course) error
	GetByID(ctx context.Context, courseID int64) (*domain.Course, error)
	List(ctx context.Context, filter CourseFilter) ([]*domain.Course, error)
	Update(ctx context.Context, course *domain.Course) error
	Delete(ctx context.Context, courseID int64) error
}

type EnrollmentFilter struct {
	StudentID *int64
	CourseID  *int64
}
type EnrollmentRepository interface {
	Create(ctx context.Context, enrollment *domain.Enrollment) error
	GetByID(ctx context.Context, enrollmentID int64) (*domain.Enrollment, error)
	List(ctx context.Context, filter EnrollmentFilter) ([]*domain.Enrollment, error)
	Update(ctx context.Context, enrollment *domain.Enrollment) error // TODO: Add status editing later...
	Delete(ctx context.Context, enrollmentID int64) error
}

type TestFilter struct {
	CourseID *int64
	Name     *string
	TestDate *time.Time
}
type TestRepository interface {
	Create(ctx context.Context, test *domain.Test) error
	GetByID(ctx context.Context, testID int64) (*domain.Test, error)
	List(ctx context.Context, filter TestFilter) ([]*domain.Test, error)
	Update(ctx context.Context, test *domain.Test) error
	Delete(ctx context.Context, testID int64) error
}

type GradeFilter struct {
	TestID       *int64
	EnrollmentID *int64
}
type GradeRepository interface {
	Create(ctx context.Context, grade *domain.Grade) error
	GetByID(ctx context.Context, gradeID int64) (*domain.Grade, error)
	List(ctx context.Context, filter GradeFilter) ([]*domain.Grade, error)
	Update(ctx context.Context, grade *domain.Grade) error
	Delete(ctx context.Context, gradeID int64) error
}
