package http

import (
	"log/slog"

	"github.com/alvarolucio2007/Scholarly/internal/services"
)

type Handler struct {
	users       *services.UserService
	students    *services.StudentService
	teachers    *services.TeacherService
	courses     *services.CourseService
	enrollments *services.EnrollmentService
	tests       *services.TestService
	grades      *services.GradeService
	logger      *slog.Logger
}

func NewHandler(
	users *services.UserService,
	students *services.StudentService,
	teachers *services.TeacherService,
	courses *services.CourseService,
	enrollments *services.EnrollmentService,
	tests *services.TestService,
	grades *services.GradeService,
	logger *slog.Logger,
) *Handler {
	return &Handler{
		users:       users,
		students:    students,
		teachers:    teachers,
		courses:     courses,
		enrollments: enrollments,
		tests:       tests,
		grades:      grades,
		logger:      logger,
	}
}

func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func nilIfZeroInt(i int64) *int64 {
	if i == 0 {
		return nil
	}
	return &i
}

func nilIfZeroFloat(f float64) *float64 {
	if f == 0 {
		return nil
	}
	return &f
}
