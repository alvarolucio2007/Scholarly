package services

import "github.com/alvarolucio2007/Scholarly/internal/ports"

type UserService struct {
	users  ports.UserRepository
	hasher ports.PasswordHasher
}
