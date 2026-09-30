package services

import (
	"context"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/alvarolucio2007/Scholarly/internal/ports"
)

type CourseService struct {
	courses ports.CourseRepository
}

func NewCourseService(courses ports.CourseRepository) *CourseService {
	return &CourseService{courses: courses}
}

type CreateCoursePayload struct {
	TeacherID   int64
	Name        string
	Code        string
	Semester    string
	MaxStudents int16
}

func (s *CourseService) CreateCourse(ctx context.Context, payload CreateCoursePayload) (*domain.Course, error) {
	course := &domain.Course{TeacherID: payload.TeacherID, Name: payload.Name, Code: payload.Code, Semester: payload.Semester, MaxStudents: payload.MaxStudents}
	if err := s.courses.Create(ctx, course); err != nil {
		return nil, err
	}
	return course, nil
}

func (s *CourseService) GetByID(ctx context.Context, courseID int64) (*domain.Course, error) {
	return s.courses.GetByID(ctx, courseID)
}

func (s *CourseService) ListCourses(ctx context.Context, filter ports.CourseFilter) ([]*domain.Course, error) {
	return s.courses.List(ctx, filter)
}

type UpdateCoursePayload struct {
	ID          int64
	TeacherID   *int64
	Name        *string
	Code        *string
	Semester    *string
	MaxStudents *int16
}

func (s *CourseService) UpdateCourse(ctx context.Context, payload UpdateCoursePayload) (*domain.Course, error) {
	course := domain.Course{ID: payload.ID}
	if payload.TeacherID != nil {
		course.TeacherID = *payload.TeacherID
	}
	if payload.Name != nil {
		course.Name = *payload.Name
	}
	if payload.Code != nil {
		course.Code = *payload.Code
	}
	if payload.Semester != nil {
		course.Semester = *payload.Semester
	}
	if payload.MaxStudents != nil {
		course.MaxStudents = *payload.MaxStudents
	}
	if err := s.courses.Update(ctx, &course); err != nil {
		return nil, err
	}
	return &course, nil
}

func (s *CourseService) DeleteCourse(ctx context.Context, courseID int64) error {
	return s.courses.Delete(ctx, courseID)
}

func (s *CourseService) EnrollStudent(ctx context.Context, studentID, courseID int64) (int64, error) {
	return s.courses.EnrollStudent(ctx, studentID, courseID)
}
