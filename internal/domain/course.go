package domain

import "time"

type Course struct {
	ID          int64
	TeacherID   int64
	Name        string
	Code        string
	Semester    string
	MaxStudents int16
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}
