package domain

import "errors"

var (
	// User
	ErrInvalidCPF         = errors.New("Invalid CPF")
	ErrInvalidEmail       = errors.New("Invalid email")
	ErrInvalidCredentials = errors.New("Invalid credentials")

	// Student
	ErrStudentNotFound          = errors.New("Student not found")
	ErrEnrollmentNumberRequired = errors.New("Enrollment number is required")

	// Teacher
	ErrDepartmentRequired = errors.New("Department is required")

	// Course
	ErrCourseNotFound = errors.New("Course not found")

	// Enrollment
	ErrAlreadyEnrolled = errors.New("Student already enrolled")
	ErrCourseFull      = errors.New("Course is full")
	ErrCannotCancel    = errors.New("Enrollment cannot be cancelled")

	// Grade
	ErrInvalidGrade       = errors.New("Invalid grade value")
	ErrGradeAlreadyExists = errors.New("Grade already exists")

	// General
	ErrNotFound                = errors.New("Not found")
	ErrEmailAlreadyExists      = errors.New("Email already exists")
	ErrCPFAlreadyExists        = errors.New("CPF already exists")
	ErrEnrollmentAlreadyExists = errors.New("Enrollment number already exists")
	ErrCourseCodeAlreadyExists = errors.New("Course code already exists for this semester")
	ErrConflict                = errors.New("Conflict")
	ErrForeignKeyViolation     = errors.New("Referenced entity does not exist")
)
