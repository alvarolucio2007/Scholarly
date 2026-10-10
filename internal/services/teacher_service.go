package services

import (
	"context"
	"log/slog"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/alvarolucio2007/Scholarly/internal/ports"
)

type TeacherService struct {
	teachers ports.TeacherRepository
	cache    ports.Cache[domain.Teacher]
	logger   *slog.Logger
}

func NewTeacherService(teachers ports.TeacherRepository, cache ports.Cache[domain.Teacher], logger *slog.Logger) *TeacherService {
	return &TeacherService{teachers: teachers, cache: cache, logger: logger}
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
	if err := s.cache.Create(ctx, teacher); err != nil {
		s.logger.ErrorContext(ctx, "failed to cache teacher", "error", err)
	}
	return teacher, nil
}

func (s *TeacherService) GetByID(ctx context.Context, teacherID int64) (*domain.Teacher, error) {
	teacher, err := s.cache.Read(ctx, teacherID)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to fetch cached teacher", "error", err)
	}
	if teacher != nil {
		return teacher, nil
	}
	teacher, err = s.teachers.GetByID(ctx, teacherID)
	if err != nil {
		return nil, err
	}
	if err := s.cache.Create(ctx, teacher); err != nil {
		s.logger.ErrorContext(ctx, "failed to cache teacher", "error", err)
	}
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
	if err := s.cache.Delete(ctx, teacher.UserID); err != nil {
		s.logger.ErrorContext(ctx, "failed to delete cached teacher", "error", err)
	}
	if err := s.cache.Create(ctx, &teacher); err != nil {
		s.logger.ErrorContext(ctx, "failed to cache teacher", "error", err)
	}
	return &teacher, nil
}

func (s *TeacherService) Delete(ctx context.Context, teacherID int64) error {
	if err := s.teachers.Delete(ctx, teacherID); err != nil {
		return err
	}
	if err := s.cache.Delete(ctx, teacherID); err != nil {
		s.logger.ErrorContext(ctx, "failed to delete cached teacher", "error", err)
	}
	return nil
}
