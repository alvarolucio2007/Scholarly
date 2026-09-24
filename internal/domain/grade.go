package domain

import "time"

type Grade struct {
	ID           int64
	TestID       int64
	EnrollmentID int64
	Value        float64
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}
