package domain

import (
	"errors"
	"strings"
)

var ErrEnrollmentNumberRequired = errors.New("enrollment number is required")

type Student struct {
	UserID           int64
	EnrollmentNumber string
}

func (s *Student) Validate() error {
	if strings.TrimSpace(s.EnrollmentNumber) == "" {
		return ErrEnrollmentNumberRequired
	}
	return nil
}
