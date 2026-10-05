package main

import (
	"log"
	"os"
	"time"

	"github.com/alvarolucio2007/Scholarly/internal/adapters/argon2"
	"github.com/alvarolucio2007/Scholarly/internal/adapters/http"
	"github.com/alvarolucio2007/Scholarly/internal/adapters/postgres"
	"github.com/alvarolucio2007/Scholarly/internal/services"
)

func main() {
	DBUrl := os.Getenv("DATABASE_URL")
	if DBUrl == "" {
		DBUrl = "postgres://admin:root@localhost:5432/scholarly?sslmode=disable"
	}
	log.Printf("starting app...")
	log.Printf("trying to connect to DB")
	db, err := postgres.NewConn(DBUrl, 30, 30, time.Minute)
	if err != nil {
		log.Fatalf("main: error while connecting to db: %v", err)
	}
	log.Printf("DB connection successful")
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("error while closing DB: %v", err)
		}
	}()
	userRepo := postgres.NewUserRepo(db)
	studentRepo := postgres.NewStudentRepo(db)
	teacherRepo := postgres.NewTeacherRepo(db)
	courseRepo := postgres.NewCourseRepo(db)
	enrollmentRepo := postgres.NewEnrollmentRepo(db)
	testRepo := postgres.NewTestRepo(db)
	gradeRepo := postgres.NewGradeRepo(db)

	hasher := argon2.NewArgon2Hasher(64*1024, 3, 2)

	userSvc := services.NewUserService(userRepo, hasher)
	studentSvc := services.NewStudentService(studentRepo)
	teacherSvc := services.NewTeacherService(teacherRepo)
	courseSvc := services.NewCourseService(courseRepo)
	enrollmentSvc := services.NewEnrollmentService(enrollmentRepo)
	testSvc := services.NewTestService(testRepo)
	gradeSvc := services.NewGradeService(gradeRepo)

	handler := http.NewHandler(userSvc, studentSvc, teacherSvc,
		courseSvc, enrollmentSvc, testSvc, gradeSvc)
	log.Printf("trying to create API")

	addr := os.Getenv("SERVER_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8080"
	}
	if err := handler.Run(handler.Mount(), addr); err != nil {
		log.Fatalf("main: error while booting HTTP server: %v", err)
	}
	log.Printf("server is up")
}
