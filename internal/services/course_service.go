package services

import (
	"context"
	"log/slog"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/alvarolucio2007/Scholarly/internal/ports"
)

type CourseService struct {
	repo   ports.CourseRepository
	cache  ports.Cache[domain.Course]
	logger *slog.Logger
}

func NewCourseService(repo ports.CourseRepository, cache ports.Cache[domain.Course], logger *slog.Logger) *CourseService {
	return &CourseService{repo: repo, cache: cache, logger: logger}
}

type CreateCoursePayload struct {
	TeacherID   int64
	Name        string
	Code        string
	Semester    string
	MaxStudents int16
}

func (s *CourseService) Create(ctx context.Context, payload CreateCoursePayload) (domain.Course, error) {
	course := domain.Course{TeacherID: payload.TeacherID, Name: payload.Name, Code: payload.Code, Semester: payload.Semester, MaxStudents: payload.MaxStudents}
	if err := s.repo.Create(ctx, &course); err != nil {
		return domain.Course{}, err
	}
	if err := s.cache.Create(ctx, course); err != nil {
		s.logger.ErrorContext(ctx, "failed to cache course", "error", err)
	}
	return course, nil
}

func (s *CourseService) GetByID(ctx context.Context, courseID int64) (domain.Course, error) {
	course, err := s.cache.Read(ctx, courseID)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to fetch cached course", "error", err)
	}
	if course.ID != 0 {
		return course, nil
	}
	course, err = s.repo.GetByID(ctx, courseID)
	if err != nil {
		return domain.Course{}, err
	}
	if err := s.cache.Create(ctx, course); err != nil {
		s.logger.ErrorContext(ctx, "failed to cache user", "error", err)
	}
	return course, nil
}

func (s *CourseService) List(ctx context.Context, filter ports.CourseFilter) ([]domain.Course, error) {
	return s.repo.List(ctx, filter)
}

type UpdateCoursePayload struct {
	ID          int64
	TeacherID   *int64
	Name        *string
	Code        *string
	Semester    *string
	MaxStudents *int16
}

func (s *CourseService) Update(ctx context.Context, payload UpdateCoursePayload) error {
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
	if err := s.repo.Update(ctx, course); err != nil {
		return err
	}
	if err := s.cache.Delete(ctx, course.ID); err != nil {
		s.logger.ErrorContext(ctx, "failed to delete cached course", "error", err)
	}
	return nil
}

func (s *CourseService) Delete(ctx context.Context, courseID int64) error {
	if err := s.repo.Delete(ctx, courseID); err != nil {
		return err
	}
	if err := s.cache.Delete(ctx, courseID); err != nil {
		s.logger.ErrorContext(ctx, "failed to delete cached user", "error", err)
	}
	return nil
}

func (s *CourseService) EnrollStudent(ctx context.Context, studentID, courseID int64) (int64, error) {
	return s.repo.EnrollStudent(ctx, studentID, courseID)
}
