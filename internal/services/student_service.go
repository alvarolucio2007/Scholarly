package services

import (
	"context"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/alvarolucio2007/Scholarly/internal/ports"
)

type StudentService struct {
	students ports.StudentRepository
}

func NewStudentService(students ports.StudentRepository) *StudentService {
	return &StudentService{students: students}
}

type CreateStudentPayload struct {
	UserID           int64
	EnrollmentNumber string
}

func (s *StudentService) CreateStudent(ctx context.Context, payload CreateStudentPayload) error {
	student := &domain.Student{UserID: payload.UserID, EnrollmentNumber: payload.EnrollmentNumber}
	return s.students.Create(ctx, student)
}

func (s *StudentService) GetByID(ctx context.Context, studentID int64) (*domain.Student, error) {
	return s.students.GetByID(ctx, studentID)
}

func (s *StudentService) GetByEnrollment(ctx context.Context, enrollment string) (*domain.Student, error) {
	return s.students.GetByEnrollment(ctx, enrollment)
}

type UpdateStudentPayload struct {
	UserID           *int64
	EnrollmentNumber *string
}

func (s *StudentService) Update(ctx context.Context, payload UpdateStudentPayload) error {
	var student domain.Student
	if payload.UserID != nil {
		student.UserID = *payload.UserID
	}
	if payload.EnrollmentNumber != nil {
		student.EnrollmentNumber = *payload.EnrollmentNumber
	}
	if err := s.students.Update(ctx, &student); err != nil {
		return err
	}
	return nil
}

func (s *StudentService) Delete(ctx context.Context, studentID int64) error {
	return s.students.Delete(ctx, studentID)
}
