package http

import (
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
}

func NewHandler(
	users *services.UserService,
	students *services.StudentService,
	teachers *services.TeacherService,
	courses *services.CourseService,
	enrollments *services.EnrollmentService,
	tests *services.TestService,
	grades *services.GradeService,
) *Handler {
	return &Handler{
		users:       users,
		students:    students,
		teachers:    teachers,
		courses:     courses,
		enrollments: enrollments,
		tests:       tests,
		grades:      grades,
	}
}

func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
