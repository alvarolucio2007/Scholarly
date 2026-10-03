package services

import (
	"context"
	"fmt"
	"time"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/alvarolucio2007/Scholarly/internal/ports"
)

type UserService struct {
	users  ports.UserRepository
	hasher ports.PasswordHasher
}

func NewUserService(users ports.UserRepository, hasher ports.PasswordHasher) *UserService {
	return &UserService{users: users, hasher: hasher}
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

func (s *UserService) CreateUser(ctx context.Context, payload CreateUserPayload) (*UserResponse, error) {
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

	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}
	return &UserResponse{ID: user.ID, Name: user.Name, CPF: user.CPF, Email: user.Email, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt}, nil
}

func (s *UserService) GetUserByID(ctx context.Context, userID int64) (*domain.User, error) {
	return s.users.GetByID(ctx, userID)
}

func (s *UserService) ListUsers(ctx context.Context, filter ports.UserFilter) ([]*domain.User, error) {
	return s.users.List(ctx, filter)
}

type UpdateUserPayload struct {
	ID       int64
	Name     *string
	CPF      *string
	Email    *string
	Password *string
}

func (s *UserService) UpdateUser(ctx context.Context, payload UpdateUserPayload) (*domain.User, error) {
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
		user.Email = *payload.Email
	}
	if payload.Password != nil {
		hash, err := s.hasher.Hash(*payload.Password)
		if err != nil {
			return nil, fmt.Errorf("service: hash password: %w", err)
		}
		user.PasswordHash = []byte(hash)
	}

	if err := s.users.Update(ctx, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) Delete(ctx context.Context, userID int64) error {
	return s.users.Delete(ctx, userID)
}
