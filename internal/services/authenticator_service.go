package services

import (
	"github.com/alvarolucio2007/Scholarly/internal/ports"
)

type AuthenticatorService struct {
	auth ports.Authenticator
}

func (s *AuthenticatorService) CreateToken(claims ports.Claims) (string, error) {
	return s.auth.CreateToken(claims)
}
func (s *AuthenticatorService) VerifyToken(tokenString string) (ports.Claims, error) {
	return s.auth.VerifyToken(tokenString)
}
