package domain

import "errors"

var (
	// User
	ErrInvalidCPF = errors.New("invalid CPF")

	// Enrollment
	ErrAlreadyEnrolled = errors.New("student already enrolled")
	ErrCourseFull      = errors.New("course is full")
	ErrCannotCancel    = errors.New("enrollment cannot be cancelled")

	// Teacher
	ErrDepartmentRequired = errors.New("department is required")

	// Grade
	ErrInvalidGrade = errors.New("invalid grade value")

	// General
	ErrNotFound = errors.New("not found")
)
