package services

import (
	"context"
	"log/slog"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/alvarolucio2007/Scholarly/internal/ports"
)

type EnrollmentService struct {
	repo   ports.EnrollmentRepository
	cache  ports.Cache[domain.Enrollment]
	logger *slog.Logger
}

func NewEnrollmentService(repo ports.EnrollmentRepository, cache ports.Cache[domain.Enrollment], logger *slog.Logger) *EnrollmentService {
	return &EnrollmentService{repo: repo, cache: cache, logger: logger}
}

type CreateEnrollmentPayload struct {
	StudentID int64
	CourseID  int64
	Status    domain.EnrollmentStatus
}

func (s *EnrollmentService) Create(ctx context.Context, payload CreateEnrollmentPayload) (*domain.Enrollment, error) {
	enrollment := &domain.Enrollment{StudentID: payload.StudentID, CourseID: payload.CourseID, Status: payload.Status}
	if err := s.repo.Create(ctx, enrollment); err != nil {
		return nil, err
	}
	return enrollment, nil
}

func (s *EnrollmentService) GetByID(ctx context.Context, enrollmentID int64) (*domain.Enrollment, error) {
	enrollment, err := s.cache.Read(ctx, enrollmentID)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to fetch cached enrollment", "error", err)
	}
	if enrollment != nil {
		return enrollment, nil
	}
	enrollment, err = s.repo.GetByID(ctx, enrollmentID)
	if err != nil {
		return nil, err
	}
	if err := s.cache.Create(ctx, enrollment); err != nil {
		s.logger.ErrorContext(ctx, "failed to cache enrollment", "error", err)
	}
	return enrollment, nil
}

func (s *EnrollmentService) List(ctx context.Context, filter ports.EnrollmentFilter) ([]*domain.Enrollment, error) {
	return s.repo.List(ctx, filter)
}

type UpdateEnrollmentPayload struct {
	ID        int64
	StudentID *int64
	CourseID  *int64
}

func (s *EnrollmentService) Update(ctx context.Context, payload UpdateEnrollmentPayload) (*domain.Enrollment, error) {
	enrollment := domain.Enrollment{ID: payload.ID}
	if payload.StudentID != nil {
		enrollment.StudentID = *payload.StudentID
	}
	if payload.CourseID != nil {
		enrollment.CourseID = *payload.CourseID
	}
	if err := s.repo.Update(ctx, &enrollment); err != nil {
		return nil, err
	}
	if err := s.cache.Delete(ctx, enrollment.ID); err != nil {
		s.logger.ErrorContext(ctx, "failed to delete cached user", "error", err)
	}
	if err := s.cache.Create(ctx, &enrollment); err != nil {
		s.logger.ErrorContext(ctx, "failed to cache user", "error", err)
	}
	return &enrollment, nil
}

func (s *EnrollmentService) Delete(ctx context.Context, enrollmentID int64) error {
	if err := s.repo.Delete(ctx, enrollmentID); err != nil {
		return err
	}
	if err := s.cache.Delete(ctx, enrollmentID); err != nil {
		s.logger.ErrorContext(ctx, "failed to delete cached enrollment", "error", err)
	}
	return nil
}
