package services

import (
	"context"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/alvarolucio2007/Scholarly/internal/ports"
)

type GradeService struct {
	grades ports.GradeRepository
}

func NewGradeService(grades ports.GradeRepository) *GradeService {
	return &GradeService{grades: grades}
}

type CreateGradePayload struct {
	TestID       int64
	EnrollmentID int64
	Value        float64
}

func (s *GradeService) CreateGrade(ctx context.Context, payload CreateGradePayload) (*domain.Grade, error) {
	grade := &domain.Grade{TestID: payload.TestID, EnrollmentID: payload.EnrollmentID, Value: payload.Value}
	if err := s.grades.Create(ctx, grade); err != nil {
		return nil, err
	}
	return grade, nil
}

func (s *GradeService) GetByID(ctx context.Context, gradeID int64) (*domain.Grade, error) {
	return s.grades.GetByID(ctx, gradeID)
}

func (s *GradeService) List(ctx context.Context, filter ports.GradeFilter) ([]*domain.Grade, error) {
	return s.grades.List(ctx, filter)
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
	if err := s.grades.Update(ctx, &grade); err != nil {
		return nil, err
	}
	return &grade, nil
}

func (s *GradeService) Delete(ctx context.Context, gradeID int64) error {
	return s.grades.Delete(ctx, gradeID)
}

func (s *GradeService) ListAverages(ctx context.Context, filter ports.ReportCardFilter) ([]*domain.ReportCard, error) {
	return s.grades.ListAverages(ctx, filter)
}
