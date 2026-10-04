package services

import (
	"context"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/alvarolucio2007/Scholarly/internal/ports"
)

type TeacherService struct {
	teachers ports.TeacherRepository
}

func NewTeacherService(teachers ports.TeacherRepository) *TeacherService {
	return &TeacherService{teachers: teachers}
}

type CreateTeacherPayload struct {
	UserID     int64
	Department string
}

func (s *TeacherService) CreateTeacher(ctx context.Context, payload CreateTeacherPayload) (*domain.Teacher, error) {
	teacher := &domain.Teacher{UserID: payload.UserID, Department: payload.Department}
	if err := s.teachers.Create(ctx, teacher); err != nil {
		return nil, err
	}
	return teacher, nil
}

func (s *TeacherService) GetByID(ctx context.Context, teacherID int64) (*domain.Teacher, error) {
	return s.teachers.GetByID(ctx, teacherID)
}

func (s *TeacherService) List(ctx context.Context, filter ports.TeacherFilter) ([]*domain.Teacher, error) {
	return s.teachers.List(ctx, filter)
}

type UpdateTeacherPayload struct {
	UserID     int64
	Department *string
}

func (s *TeacherService) Update(ctx context.Context, payload UpdateTeacherPayload) (*domain.Teacher, error) {
	teacher := domain.Teacher{UserID: payload.UserID}
	if payload.Department != nil {
		teacher.Department = *payload.Department
	}
	if err := s.teachers.Update(ctx, &teacher); err != nil {
		return nil, err
	}
	return &teacher, nil
}

func (s *TeacherService) Delete(ctx context.Context, teacherID int64) error {
	return s.teachers.Delete(ctx, teacherID)
}
