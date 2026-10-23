package paseto

import (
	"crypto/rand"
	"errors"
	"time"
	"uuid"

	"github.com/alvarolucio2007/Scholarly/internal/ports"
)

var (
	ErrExpiredToken = errors.New("token is expired")
	ErrInvalidToken = errors.New("token is invalid")
)

func NewClaims(email string, duration time.Duration, userID int64) (ports.Claims, error) {
	tokenID := uuid.New()
	now := time.Now()
	payload := ports.Claims{
		ID:        tokenID,
		Email:     email,
		UserID:    userID,
		IssuedAt:  now,
		ExpiresAt: now.Add(duration),
	}
	return payload, nil
}

func isValidClaims(claims *ports.Claims) error {
	if time.Now().After(claims.ExpiresAt) {
		return ErrExpiredToken
	}
	return nil
}

func generateSymmetricKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	return key, err
}
