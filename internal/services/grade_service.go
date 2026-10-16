package services

import (
	"context"
	"log/slog"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/alvarolucio2007/Scholarly/internal/ports"
)

type GradeService struct {
	repo   ports.GradeRepository
	cache  ports.Cache[domain.Grade]
	logger *slog.Logger
}

func NewGradeService(repo ports.GradeRepository, cache ports.Cache[domain.Grade], logger *slog.Logger) *GradeService {
	return &GradeService{repo: repo, cache: cache, logger: logger}
}

type CreateGradePayload struct {
	TestID       int64
	EnrollmentID int64
	Value        float64
}

func (s *GradeService) Create(ctx context.Context, payload CreateGradePayload) (*domain.Grade, error) {
	grade := &domain.Grade{TestID: payload.TestID, EnrollmentID: payload.EnrollmentID, Value: payload.Value}
	if err := s.repo.Create(ctx, grade); err != nil {
		return nil, err
	}
	if err := s.cache.Create(ctx, grade); err != nil {
		s.logger.ErrorContext(ctx, "failed to cache grade", "error", err)
	}
	return grade, nil
}

func (s *GradeService) GetByID(ctx context.Context, gradeID int64) (*domain.Grade, error) {
	grade, err := s.cache.Read(ctx, gradeID)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to fetch cached grade", "error", err)
	}
	if grade != nil {
		return grade, nil
	}
	grade, err = s.repo.GetByID(ctx, gradeID)
	if err != nil {
		return nil, err
	}
	if err := s.cache.Create(ctx, grade); err != nil {
		s.logger.ErrorContext(ctx, "failed to cache grade", "error", err)
	}
	return grade, nil
}

func (s *GradeService) List(ctx context.Context, filter ports.GradeFilter) ([]*domain.Grade, error) {
	return s.repo.List(ctx, filter)
}

type UpdateGradePayload struct {
	ID           int64
	TestID       *int64
	EnrollmentID *int64
	Value        *float64
}

func (s *GradeService) Update(ctx context.Context, payload UpdateGradePayload) (*domain.Grade, error) {
	grade := domain.Grade{ID: payload.ID}
	if payload.TestID != nil {
		grade.TestID = *payload.TestID
	}
	if payload.EnrollmentID != nil {
		grade.EnrollmentID = *payload.EnrollmentID
	}
	if payload.Value != nil {
		grade.Value = *payload.Value
	}
	if err := s.repo.Update(ctx, &grade); err != nil {
		return nil, err
	}
	if err := s.cache.Delete(ctx, grade.ID); err != nil {
		s.logger.ErrorContext(ctx, "failed to delete cached grade", "error", err)
	}
	if err := s.cache.Create(ctx, &grade); err != nil {
		s.logger.ErrorContext(ctx, "failed to cache grade", "error", err)
	}
	return &grade, nil
}

func (s *GradeService) Delete(ctx context.Context, gradeID int64) error {
	if err := s.repo.Delete(ctx, gradeID); err != nil {
		return err
	}
	if err := s.cache.Delete(ctx, gradeID); err != nil {
		s.logger.ErrorContext(ctx, "failed to delete cached grade", "error", err)
	}
	return nil
}

func (s *GradeService) ListAverages(ctx context.Context, filter ports.ReportCardFilter) ([]*domain.ReportCard, error) {
	return s.repo.ListAverages(ctx, filter)
}

func (s *GradeService) GetAverage(ctx context.Context, studentID, courseID int64) (float64, error) {
	return s.repo.GetAverage(ctx, studentID, courseID)
}
