package domain

import "time"

type Tests struct {
	ID        int64
	CourseID  int64
	Name      string
	Weight    float64
	TestDate  time.Time
	CreatedAt time.Time
	UpdatedAt *time.Time
}
