package paseto

import (
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
