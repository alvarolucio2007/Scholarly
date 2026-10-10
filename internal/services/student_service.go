package services

import (
	"context"
	"log/slog"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/alvarolucio2007/Scholarly/internal/ports"
)

type StudentService struct {
	students ports.StudentRepository
	cache    ports.Cache[domain.Student]
	logger   *slog.Logger
}

func NewStudentService(students ports.StudentRepository, cache ports.Cache[domain.Student], logger *slog.Logger) *StudentService {
	return &StudentService{students: students, cache: cache, logger: logger}
}

type CreateStudentPayload struct {
	UserID           int64
	EnrollmentNumber string
}

func (s *StudentService) CreateStudent(ctx context.Context, payload CreateStudentPayload) (*domain.Student, error) {
	student := &domain.Student{UserID: payload.UserID, EnrollmentNumber: payload.EnrollmentNumber}
	if err := s.students.Create(ctx, student); err != nil {
		return nil, err
	}
	if err := s.cache.Create(ctx, student); err != nil {
		s.logger.ErrorContext(ctx, "failed to cache student", "error", err)
	}
	return student, nil
}

func (s *StudentService) GetByID(ctx context.Context, studentID int64) (*domain.Student, error) {
	student, err := s.cache.Read(ctx, studentID)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to fetch cached student", "error", err)
	}
	if student != nil {
		return student, nil
	}
	student, err = s.students.GetByID(ctx, studentID)
	if err != nil {
		return nil, err
	}
	if err := s.cache.Create(ctx, student); err != nil {
		s.logger.ErrorContext(ctx, "failed to cache student", "error", err)
	}
	return student, nil
}

func (s *StudentService) GetByEnrollment(ctx context.Context, enrollment string) (*domain.Student, error) {
	return s.students.GetByEnrollment(ctx, enrollment)
}

type UpdateStudentPayload struct {
	UserID           int64
	EnrollmentNumber *string
}

func (s *StudentService) Update(ctx context.Context, payload UpdateStudentPayload) (*domain.Student, error) {
	student := domain.Student{UserID: payload.UserID}
	if payload.EnrollmentNumber != nil {
		student.EnrollmentNumber = *payload.EnrollmentNumber
	}
	if err := s.students.Update(ctx, &student); err != nil {
		return nil, err
	}
	if err := s.cache.Delete(ctx, student.UserID); err != nil {
		s.logger.ErrorContext(ctx, "failed to delete cached student", "error", err)
	}
	if err := s.cache.Create(ctx, &student); err != nil {
		s.logger.ErrorContext(ctx, "failed to cache student", "error", err)
	}
	return &student, nil
}

func (s *StudentService) Delete(ctx context.Context, studentID int64) error {
	if err := s.students.Delete(ctx, studentID); err != nil {
		return err
	}
	if err := s.cache.Delete(ctx, studentID); err != nil {
		s.logger.ErrorContext(ctx, "failed to delete cached student", "error", err)
	}
	return nil
}
