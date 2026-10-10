package services

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/alvarolucio2007/Scholarly/internal/ports"
)

type UserService struct {
	repo   ports.UserRepository
	hasher ports.PasswordHasher
	cache  ports.Cache[domain.User]
	logger *slog.Logger
}

func NewUserService(repo ports.UserRepository, hasher ports.PasswordHasher, cache ports.Cache[domain.User], logger *slog.Logger) *UserService {
	return &UserService{repo: repo, hasher: hasher, cache: cache, logger: logger}
}

type CreateUserPayload struct {
	Name     string
	CPF      string
	Email    string
	Password string
}
type UserResponse struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	CPF       string     `json:"cpf"`
	Email     string     `json:"email"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func (s *UserService) Create(ctx context.Context, payload CreateUserPayload) (*domain.User, error) {
	user := &domain.User{
		Name:  payload.Name,
		CPF:   payload.CPF,
		Email: payload.Email,
	}
	if err := user.ValidateCPF(); err != nil {
		return nil, err
	}

	hash, err := s.hasher.Hash(payload.Password)
	if err != nil {
		return nil, fmt.Errorf("service: hash password: %w", err)
	}
	user.PasswordHash = []byte(hash)

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	if err := s.cache.Create(ctx, user); err != nil {
		s.logger.ErrorContext(ctx, "failed to cache user", "error", err)
	}
	return user, nil
}

func (s *UserService) GetByID(ctx context.Context, userID int64) (*domain.User, error) {
	user, err := s.cache.Read(ctx, userID)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to fetch cached user", "error", err)
	}
	if user != nil {
		return user, nil
	}
	user, err = s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if err := s.cache.Create(ctx, user); err != nil {
		s.logger.ErrorContext(ctx, "failed to cache user", "error", err)
	}
	return user, nil
}

func (s *UserService) List(ctx context.Context, filter ports.UserFilter) ([]*domain.User, error) {
	return s.repo.List(ctx, filter)
}

type UpdateUserPayload struct {
	ID       int64
	Name     *string
	CPF      *string
	Email    *string
	Password *string
}

func (s *UserService) Update(ctx context.Context, payload UpdateUserPayload) (*domain.User, error) {
	user := domain.User{ID: payload.ID}
	if payload.Name != nil {
		user.Name = *payload.Name
	}
	if payload.CPF != nil {
		user.CPF = *payload.CPF
		if err := user.ValidateCPF(); err != nil {
			return nil, err
		}
	}
	if payload.Email != nil {
		if err := user.ValidateEmail(); err != nil {
			return nil, err
		}
		user.Email = *payload.Email
	}
	if payload.Password != nil {
		hash, err := s.hasher.Hash(*payload.Password)
		if err != nil {
			return nil, fmt.Errorf("service: hash password: %w", err)
		}
		user.PasswordHash = []byte(hash)
	}

	if err := s.repo.Update(ctx, &user); err != nil {
		return nil, err
	}
	if err := s.cache.Delete(ctx, user.ID); err != nil {
		s.logger.ErrorContext(ctx, "failed to delete cached user", "error", err)
	}
	if err := s.cache.Create(ctx, &user); err != nil {
		s.logger.ErrorContext(ctx, "failed to cache user", "error", err)
	}
	return &user, nil
}

func (s *UserService) Delete(ctx context.Context, userID int64) error {
	if err := s.repo.Delete(ctx, userID); err != nil {
		return err
	}
	if err := s.cache.Delete(ctx, userID); err != nil {
		s.logger.ErrorContext(ctx, "failed to delete cached user", "error", err)
	}
	return nil
}
