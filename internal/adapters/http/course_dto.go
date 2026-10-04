package http

import "time"

type CreateCourseDTO struct {
	TeacherID   int64  `json:"teacher_id" validate:"required,gte=1"`
	Name        string `json:"name" validate:"required,max=255"`
	Code        string `json:"code" validate:"required,max=50"`
	Semester    string `json:"semester" validate:"required,max=50"`
	MaxStudents int16  `json:"max_students" validate:"required,gte=1"`
}
type UpdateCourseDTO struct {
	ID          int64   `json:"id" validate:"required,gte=1"`
	TeacherID   *int64  `json:"teacher_id" validate:"gte=1"`
	Name        *string `json:"name" validate:"max=50"`
	Code        *string `json:"code" validate:"max=50"`
	Semester    *string `json:"semester" validate:"max=50"`
	MaxStudents *int16  `json:"max_students" validate:"max=50"`
}
type CourseResponseDTO struct {
	ID          int64      `json:"id"`
	TeacherID   int64      `json:"teacher_id"`
	Name        string     `json:"name"`
	Code        string     `json:"code"`
	Semester    string     `json:"semester"`
	MaxStudents int16      `json:"max_students"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}
