package ports

import (
	"context"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
)

type UserCache interface {
	Create(ctx context.Context, user *domain.User) error
	Read(ctx context.Context, userID int64) (*domain.User, error)
	Delete(ctx context.Context, userID int64) error
}
type TeacherCache interface {
	Create(ctx context.Context, teacher *domain.Teacher) error
	Read(ctx context.Context, teacherID int64) (*domain.Teacher, error)
	Delete(ctx context.Context, teacherID int64) error
}
type StudentCache interface {
	Create(ctx context.Context, student *domain.Student) error
	Read(ctx context.Context, studentID int64) (*domain.Student, error)
	Delete(ctx context.Context, studentID int64) error
}
type CourseCache interface {
	Create(ctx context.Context, course *domain.Course) error
	Read(ctx context.Context, courseID int64) (*domain.Course, error)
	Delete(ctx context.Context, courseID int64) error
}
type EnrollmentCache interface {
	Create(ctx context.Context, enrollment *domain.Enrollment) error
	Read(ctx context.Context, enrollmentID int64) (*domain.Enrollment, error)
	Delete(ctx context.Context, enrollmentID int64) error
}
type TestCache interface {
	Create(ctx context.Context, test *domain.Test) error
	Read(ctx context.Context, testID int64) (*domain.Test, error)
	Delete(ctx context.Context, testID int64) error
}
type GradeCache interface {
	Create(ctx context.Context, grade *domain.Grade) error
	Read(ctx context.Context, gradeID int64) (*domain.Grade, error)
	Delete(ctx context.Context, gradeID int64) error
}
