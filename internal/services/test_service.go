package services

import (
	"context"
	"time"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/alvarolucio2007/Scholarly/internal/ports"
)

type TestService struct {
	tests ports.TestRepository
}

func NewTestService(tests ports.TestRepository) *TestService {
	return &TestService{tests: tests}
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
	return test, nil
}

func (s *TestService) GetByID(ctx context.Context, testID int64) (*domain.Test, error) {
	return s.tests.GetByID(ctx, testID)
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
	return &test, nil
}

func (s *TestService) DeleteTest(ctx context.Context, testID int64) error {
	return s.tests.Delete(ctx, testID)
}
