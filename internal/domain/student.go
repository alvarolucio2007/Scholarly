package domain

import (
	"errors"
)

var ErrEnrollmentNumberRequired = errors.New("enrollment number is required")

type Student struct {
	UserID           int64
	EnrollmentNumber string
}
