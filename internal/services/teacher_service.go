package services

import (
	"context"
	"log/slog"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/alvarolucio2007/Scholarly/internal/ports"
)

type TeacherService struct {
	repo   ports.TeacherRepository
	cache  ports.Cache[domain.Teacher]
	logger *slog.Logger
}

func NewTeacherService(repo ports.TeacherRepository, cache ports.Cache[domain.Teacher], logger *slog.Logger) *TeacherService {
	return &TeacherService{repo: repo, cache: cache, logger: logger}
}

type CreateTeacherPayload struct {
	UserID     int64
	Department string
}

func (s *TeacherService) Create(ctx context.Context, payload CreateTeacherPayload) (domain.Teacher, error) {
	teacher := domain.Teacher{UserID: payload.UserID, Department: payload.Department}
	if err := s.repo.Create(ctx, &teacher); err != nil {
		return domain.Teacher{}, err
	}
	if err := s.cache.Create(ctx, teacher); err != nil {
		s.logger.ErrorContext(ctx, "failed to cache teacher", "error", err)
	}
	return teacher, nil
}

func (s *TeacherService) GetByID(ctx context.Context, teacherID int64) (domain.Teacher, error) {
	teacher, err := s.cache.Read(ctx, teacherID)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to fetch cached teacher", "error", err)
	}
	if teacher.UserID != 0 {
		return teacher, nil
	}
	teacher, err = s.repo.GetByID(ctx, teacherID)
	if err != nil {
		return domain.Teacher{}, err
	}
	if err := s.cache.Create(ctx, teacher); err != nil {
		s.logger.ErrorContext(ctx, "failed to cache teacher", "error", err)
	}
	return s.repo.GetByID(ctx, teacherID)
}

func (s *TeacherService) List(ctx context.Context, filter ports.TeacherFilter) ([]domain.Teacher, error) {
	return s.repo.List(ctx, filter)
}

type UpdateTeacherPayload struct {
	UserID     int64
	Department *string
}

func (s *TeacherService) Update(ctx context.Context, payload UpdateTeacherPayload) error {
	teacher := domain.Teacher{UserID: payload.UserID}
	if payload.Department != nil {
		teacher.Department = *payload.Department
	}
	if err := s.repo.Update(ctx, teacher); err != nil {
		return err
	}
	if err := s.cache.Delete(ctx, teacher.UserID); err != nil {
		s.logger.ErrorContext(ctx, "failed to delete cached teacher", "error", err)
	}
	return nil
}

func (s *TeacherService) Delete(ctx context.Context, teacherID int64) error {
	if err := s.repo.Delete(ctx, teacherID); err != nil {
		return err
	}
	if err := s.cache.Delete(ctx, teacherID); err != nil {
		s.logger.ErrorContext(ctx, "failed to delete cached teacher", "error", err)
	}
	return nil
}
