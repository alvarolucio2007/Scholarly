package main

import (
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/alvarolucio2007/Scholarly/internal/adapters/argon2"
	"github.com/alvarolucio2007/Scholarly/internal/adapters/http"
	"github.com/alvarolucio2007/Scholarly/internal/adapters/logger"
	"github.com/alvarolucio2007/Scholarly/internal/adapters/postgres"
	"github.com/alvarolucio2007/Scholarly/internal/services"
)

func main() {
	logger := logger.New()
	DBUrl := os.Getenv("DATABASE_URL")
	if DBUrl == "" {
		DBUrl = "postgres://admin:root@localhost:5432/scholarly?sslmode=disable"
	}
	logger.Info("starting app")
	logger.Info("trying to connect to db")
	db, err := postgres.NewConn(DBUrl, 30, 30, time.Minute)
	if err != nil {
		log.Fatalf("main: error while connecting to db: %v", err)
	}
	logger.Info("db connection successful")
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("error while closing DB: %v", err)
		}
	}()
	userRepo := postgres.NewUserRepo(db)
	studentRepo := postgres.NewStudentRepo(db)
	teacherRepo := postgres.NewTeacherRepo(db, logger)
	courseRepo := postgres.NewCourseRepo(db, logger)
	enrollmentRepo := postgres.NewEnrollmentRepo(db, logger)
	testRepo := postgres.NewTestRepo(db, logger)
	gradeRepo := postgres.NewGradeRepo(db, logger)

	hasher := argon2.NewArgon2Hasher(64*1024, 3, 2)

	userSvc := services.NewUserService(userRepo, hasher)
	studentSvc := services.NewStudentService(studentRepo)
	teacherSvc := services.NewTeacherService(teacherRepo)
	courseSvc := services.NewCourseService(courseRepo)
	enrollmentSvc := services.NewEnrollmentService(enrollmentRepo)
	testSvc := services.NewTestService(testRepo)
	gradeSvc := services.NewGradeService(gradeRepo)

	handler := http.NewHandler(userSvc, studentSvc, teacherSvc,
		courseSvc, enrollmentSvc, testSvc, gradeSvc, logger)

	addr := os.Getenv("SERVER_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8080"
	}
	logger.Info("trying to start server")
	if err := handler.Run(handler.Mount(), addr); err != nil {
		slog.Error("main: error while booting HTTP server: ", "", err)
	}
	logger.Info("server has been shut down")
}
