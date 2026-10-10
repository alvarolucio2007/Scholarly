package services

import (
	"context"
	"log/slog"
	"time"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/alvarolucio2007/Scholarly/internal/ports"
)

type TestService struct {
	tests  ports.TestRepository
	cache  ports.Cache[domain.Test]
	logger *slog.Logger
}

func NewTestService(tests ports.TestRepository, cache ports.Cache[domain.Test], logger *slog.Logger) *TestService {
	return &TestService{tests: tests, cache: cache, logger: logger}
}

type CreateTestPayload struct {
	CourseID int64
	Name     string
	Weight   float64
	TestDate time.Time
}

func (s *TestService) CreateTest(ctx context.Context, payload CreateTestPayload) (*domain.Test, error) {
	test := &domain.Test{CourseID: payload.CourseID, Name: payload.Name, Weight: payload.Weight, TestDate: payload.TestDate}
	if err := s.tests.Create(ctx, test); err != nil {
		return nil, err
	}
	if err := s.cache.Create(ctx, test); err != nil {
		s.logger.ErrorContext(ctx, "failed to cache test", "error", err)
	}
	return test, nil
}

func (s *TestService) GetByID(ctx context.Context, testID int64) (*domain.Test, error) {
	test, err := s.cache.Read(ctx, testID)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to fetch cached tests", "error", err)
	}
	if test != nil {
		return test, nil
	}
	test, err = s.tests.GetByID(ctx, testID)
	if err != nil {
		return nil, err
	}
	if err := s.cache.Create(ctx, test); err != nil {
		s.logger.ErrorContext(ctx, "failed to cache test", "error", err)
	}
	return test, nil
}

func (s *TestService) ListTests(ctx context.Context, filter ports.TestFilter) ([]*domain.Test, error) {
	return s.tests.List(ctx, filter)
}

type UpdateTestPayload struct {
	ID       int64
	CourseID *int64
	Name     *string
	Weight   *float64
	TestDate *time.Time
}

func (s *TestService) UpdateTest(ctx context.Context, payload UpdateTestPayload) (*domain.Test, error) {
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
	if err := s.tests.Update(ctx, &test); err != nil {
		return nil, err
	}
	if err := s.cache.Delete(ctx, test.ID); err != nil {
		s.logger.ErrorContext(ctx, "failed to delete cached test", "error", err)
	}
	if err := s.cache.Create(ctx, &test); err != nil {
		s.logger.ErrorContext(ctx, "failed to cache test", "error", err)
	}
	return &test, nil
}

func (s *TestService) DeleteTest(ctx context.Context, testID int64) error {
	if err := s.tests.Delete(ctx, testID); err != nil {
		return err
	}
	if err := s.cache.Delete(ctx, testID); err != nil {
		s.logger.ErrorContext(ctx, "failed to delete cached test", "error", err)
	}
	return nil
}
