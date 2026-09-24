package domain

import "time"

type EnrollmentStatus string

const (
	StatusActive    EnrollmentStatus = "active"
	StatusCancelled EnrollmentStatus = "cancelled"
	StatusCompleted EnrollmentStatus = "completed"
)

type Enrollment struct {
	ID         int64
	StudentID  int64
	CourseID   int64
	EnrolledAt time.Time
	Status     EnrollmentStatus
}

func (s EnrollmentStatus) IsValid() bool {
	switch s {
	case StatusActive, StatusCancelled, StatusCompleted:
		return true
	}
	return false
}
