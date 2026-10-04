package services

import (
	"context"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/alvarolucio2007/Scholarly/internal/ports"
)

type EnrollmentService struct {
	enrollments ports.EnrollmentRepository
}

func NewEnrollmentService(enrollments ports.EnrollmentRepository) *EnrollmentService {
	return &EnrollmentService{enrollments: enrollments}
}

type CreateEnrollmentPayload struct {
	StudentID int64
	CourseID  int64
	Status    domain.EnrollmentStatus
}

func (s *EnrollmentService) CreateEnrollment(ctx context.Context, payload CreateEnrollmentPayload) (*domain.Enrollment, error) {
	enrollment := &domain.Enrollment{StudentID: payload.StudentID, CourseID: payload.CourseID, Status: payload.Status}
	if err := s.enrollments.Create(ctx, enrollment); err != nil {
		return nil, err
	}
	return enrollment, nil
}

func (s *EnrollmentService) GetByID(ctx context.Context, enrollmentID int64) (*domain.Enrollment, error) {
	return s.enrollments.GetByID(ctx, enrollmentID)
}

func (s *EnrollmentService) ListEnrollments(ctx context.Context, filter ports.EnrollmentFilter) ([]*domain.Enrollment, error) {
	return s.enrollments.List(ctx, filter)
}

type UpdateEnrollmentPayload struct {
	ID        int64
	StudentID *int64
	CourseID  *int64
}

func (s *EnrollmentService) UpdateEnrollment(ctx context.Context, payload UpdateEnrollmentPayload) (*domain.Enrollment, error) {
	enrollment := domain.Enrollment{ID: payload.ID}
	if payload.StudentID != nil {
		enrollment.StudentID = *payload.StudentID
	}
	if payload.CourseID != nil {
		enrollment.CourseID = *payload.CourseID
	}
	if err := s.enrollments.Update(ctx, &enrollment); err != nil {
		return nil, err
	}
	return &enrollment, nil
}

func (s *EnrollmentService) DeleteEnrollment(ctx context.Context, enrollmentID int64) error {
	return s.enrollments.Delete(ctx, enrollmentID)
}

func (s *EnrollmentService) GetWeightedAverage(ctx context.Context, studentID, courseID int64) (float64, error) {
	return s.enrollments.GetWeightedAverage(ctx, studentID, courseID)
}
