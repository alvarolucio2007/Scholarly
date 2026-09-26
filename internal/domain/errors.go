package domain

import "errors"

var (
	// User
	ErrInvalidCPF = errors.New("invalid CPF")

	// Student
	ErrEnrollmentNumberRequired = errors.New("enrollment number is required")

	// Enrollment
	ErrAlreadyEnrolled = errors.New("student already enrolled")
	ErrCourseFull      = errors.New("course is full")
	ErrCannotCancel    = errors.New("enrollment cannot be cancelled")

	// Teacher
	ErrDepartmentRequired = errors.New("department is required")

	// Grade
	ErrInvalidGrade = errors.New("invalid grade value")

	// General
	ErrNotFound                = errors.New("not found")
	ErrEmailAlreadyExists      = errors.New("email already exists")
	ErrCPFAlreadyExists        = errors.New("CPF already exists")
	ErrEnrollmentAlreadyExists = errors.New("enrollment number already exists")
	ErrCourseCodeAlreadyExists = errors.New("course code already exists for this semester")
	ErrConflict                = errors.New("conflict")
	ErrForeignKeyViolation     = errors.New("referenced entity does not exist")
)
