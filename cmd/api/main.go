package main

import (
	"log"
	"time"

	"github.com/alvarolucio2007/Scholarly/internal/adapters/postgres"
)

const DBUrl string = "postgres://admin:root@localhost:5432/scholarly?sslmode=disable"

func main() {
	db, err := postgres.NewConn(DBUrl, 30, 30, time.Minute)
	if err != nil {
		log.Printf("main: error while connecting to db: %w", err)
	}
	userRepo := postgres.NewUserRepo(db)
	studentRepo := postgres.NewStudentRepo(db)
	teacherRepo := postgres.NewTeacherRepo(db)
	courseRepo := postgres.NewCourseRepo(db)
	enrollmentRepo := postgres.NewEnrollmentRepo(db)
	testRepo := postgres.NewTestRepo(db)
	gradeREpo := postgres.NewGradeRepo(db)
}
