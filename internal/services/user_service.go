// Package services orchestrates the workflow , between the repositories, cache, hasher,logger, etc.
package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/alvarolucio2007/Scholarly/internal/adapters/paseto"
	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/alvarolucio2007/Scholarly/internal/ports"
)

type UserService struct {
	repo   ports.UserRepository
	hasher ports.PasswordHasher
	cache  ports.Cache[domain.User]
	logger *slog.Logger
	auth   ports.Authenticator
}

func NewUserService(repo ports.UserRepository, hasher ports.PasswordHasher, cache ports.Cache[domain.User], logger *slog.Logger, auth ports.Authenticator) *UserService {
	return &UserService{repo: repo, hasher: hasher, cache: cache, logger: logger, auth: auth}
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

func (s *UserService) Create(ctx context.Context, payload CreateUserPayload) (domain.User, error) {
	user := domain.User{
		Name:  payload.Name,
		CPF:   payload.CPF,
		Email: payload.Email,
	}
	if err := user.ValidateCPF(); err != nil {
		return domain.User{}, err
	}

	hash, err := s.hasher.Hash(payload.Password)
	if err != nil {
		return domain.User{}, fmt.Errorf("service: hash password: %w", err)
	}
	user.PasswordHash = []byte(hash)

	if err := s.repo.Create(ctx, &user); err != nil {
		return domain.User{}, err
	}
	if err := s.cache.Create(ctx, user); err != nil {
		s.logger.ErrorContext(ctx, "failed to cache user", "error", err)
	}
	return user, nil
}

func (s *UserService) GetByID(ctx context.Context, userID int64) (domain.User, error) {
	user, err := s.cache.Read(ctx, userID)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to fetch cached user", "error", err)
	}
	if user.ID != 0 {
		return user, nil
	}
	user, err = s.repo.GetByID(ctx, userID)
	if err != nil {
		return domain.User{}, err
	}
	if err := s.cache.Create(ctx, user); err != nil {
		s.logger.ErrorContext(ctx, "failed to cache user", "error", err)
	}
	return user, nil
}

func (s *UserService) List(ctx context.Context, filter ports.UserFilter) ([]domain.User, error) {
	return s.repo.List(ctx, filter)
}

type UpdateUserPayload struct {
	ID       int64
	Name     *string
	CPF      *string
	Email    *string
	Password *string
}

func (s *UserService) Update(ctx context.Context, payload UpdateUserPayload) error {
	user := domain.User{ID: payload.ID}
	if payload.Name != nil {
		user.Name = *payload.Name
	}
	if payload.CPF != nil {
		user.CPF = *payload.CPF
		if err := user.ValidateCPF(); err != nil {
			return err
		}
	}
	if payload.Email != nil {
		if err := user.ValidateEmail(); err != nil {
			return err
		}
		user.Email = *payload.Email
	}
	if payload.Password != nil {
		hash, err := s.hasher.Hash(*payload.Password)
		if err != nil {
			return fmt.Errorf("service: hash password: %w", err)
		}
		user.PasswordHash = []byte(hash)
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return err
	}
	if err := s.cache.Delete(ctx, user.ID); err != nil {
		s.logger.ErrorContext(ctx, "failed to delete cached user", "error", err)
	}
	return nil
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
func (s *UserService) VerifyUser(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return "", domain.ErrInvalidCredentials
		}
		return "", err
	}
	isSame, err := s.hasher.Check(password, string(user.PasswordHash))
	if err != nil {
		return "", err
	}
	if !isSame {
		return "", domain.ErrInvalidCredentials
	}
	claims, err := paseto.NewClaims(email, 24*time.Hour, user.ID)
	if err != nil {
		return "", err
	}
	token, err := s.auth.CreateToken(claims)
	if err != nil {
		return "", err
	}
	return token, nil
}
