package services

import (
	"context"
	"log/slog"
	"time"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/alvarolucio2007/Scholarly/internal/ports"
)

type TestService struct {
	repo   ports.TestRepository
	cache  ports.Cache[domain.Test]
	logger *slog.Logger
}

func NewTestService(repo ports.TestRepository, cache ports.Cache[domain.Test], logger *slog.Logger) *TestService {
	return &TestService{repo: repo, cache: cache, logger: logger}
}

type CreateTestPayload struct {
	CourseID int64
	Name     string
	Weight   float64
	TestDate time.Time
}

func (s *TestService) Create(ctx context.Context, payload CreateTestPayload) (domain.Test, error) {
	test := domain.Test{CourseID: payload.CourseID, Name: payload.Name, Weight: payload.Weight, TestDate: payload.TestDate}
	if err := s.repo.Create(ctx, &test); err != nil {
		return domain.Test{}, err
	}
	if err := s.cache.Create(ctx, test); err != nil {
		s.logger.ErrorContext(ctx, "failed to cache test", "error", err)
	}
	return test, nil
}

func (s *TestService) GetByID(ctx context.Context, testID int64) (domain.Test, error) {
	test, err := s.cache.Read(ctx, testID)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to fetch cached tests", "error", err)
	}
	if test.ID != 0 {
		return test, nil
	}
	test, err = s.repo.GetByID(ctx, testID)
	if err != nil {
		return domain.Test{}, err
	}
	if err := s.cache.Create(ctx, test); err != nil {
		s.logger.ErrorContext(ctx, "failed to cache test", "error", err)
	}
	return test, nil
}

func (s *TestService) List(ctx context.Context, filter ports.TestFilter) ([]domain.Test, error) {
	return s.repo.List(ctx, filter)
}

type UpdateTestPayload struct {
	ID       int64
	CourseID *int64
	Name     *string
	Weight   *float64
	TestDate *time.Time
}

func (s *TestService) Update(ctx context.Context, payload UpdateTestPayload) error {
	test := domain.Test{ID: payload.ID}
	if payload.Name != nil {
		test.Name = *payload.Name
	}
	if payload.CourseID != nil {
		test.CourseID = *payload.CourseID
	}
	if payload.Weight != nil {
		test.Weight = *payload.Weight
	}
	if payload.TestDate != nil {
		test.TestDate = *payload.TestDate
	}
	if err := s.repo.Update(ctx, test); err != nil {
		return err
	}
	if err := s.cache.Delete(ctx, test.ID); err != nil {
		s.logger.ErrorContext(ctx, "failed to delete cached test", "error", err)
	}
	return nil
}

func (s *TestService) Delete(ctx context.Context, testID int64) error {
	if err := s.repo.Delete(ctx, testID); err != nil {
		return err
	}
	if err := s.cache.Delete(ctx, testID); err != nil {
		s.logger.ErrorContext(ctx, "failed to delete cached test", "error", err)
	}
	return nil
}
