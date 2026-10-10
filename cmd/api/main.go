package main

import (
	"os"
	"time"

	"github.com/alvarolucio2007/Scholarly/internal/adapters/argon2"
	"github.com/alvarolucio2007/Scholarly/internal/adapters/http"
	"github.com/alvarolucio2007/Scholarly/internal/adapters/logger"
	"github.com/alvarolucio2007/Scholarly/internal/adapters/postgres"
	"github.com/alvarolucio2007/Scholarly/internal/adapters/valkey"
	"github.com/alvarolucio2007/Scholarly/internal/services"
)

func main() {
	logger := logger.New()
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://admin:root@localhost:5432/scholarly?sslmode=disable"
	}
	logger.Info("starting app")
	logger.Info("trying to connect to db")
	db, err := postgres.NewConn(dbURL, 30, 30, time.Minute)
	if err != nil {
		logger.Error("main: error while connecting to db", "error", err)
		os.Exit(1)
	}
	logger.Info("db connection successful")
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("main: error while closing db conn", "error", err)
		}
	}()
	userRepo := postgres.NewUserRepo(db, logger)
	studentRepo := postgres.NewStudentRepo(db)
	teacherRepo := postgres.NewTeacherRepo(db, logger)
	courseRepo := postgres.NewCourseRepo(db, logger)
	enrollmentRepo := postgres.NewEnrollmentRepo(db, logger)
	testRepo := postgres.NewTestRepo(db, logger)
	gradeRepo := postgres.NewGradeRepo(db, logger)

	hasher := argon2.NewArgon2Hasher(64*1024, 3, 2)

	cacheURL := os.Getenv("CACHE_URL")
	if cacheURL == "" {
		cacheURL = "localhost:6379"
	}
	logger.Info("trying to connect to cache")
	rdb := valkey.NewConn(cacheURL, "", 0)
	if rdb == nil {
		logger.Error("main: error while connecting to cache")
		os.Exit(1)
	}
	logger.Info("cache connection successful")

	userCache := valkey.NewUserCache(rdb)
	studentCache := valkey.NewStudentCache(rdb)
	teacherCache := valkey.NewTeacherCache(rdb)
	courseCache := valkey.NewCourseCache(rdb)
	enrollmentCache := valkey.NewEnrollmentCache(rdb)
	testCache := valkey.NewTestCache(rdb)
	gradeCache := valkey.NewGradeCache(rdb)

	userSvc := services.NewUserService(userRepo, hasher, userCache, logger)
	studentSvc := services.NewStudentService(studentRepo, studentCache, logger)
	teacherSvc := services.NewTeacherService(teacherRepo, teacherCache, logger)
	courseSvc := services.NewCourseService(courseRepo, courseCache, logger)
	enrollmentSvc := services.NewEnrollmentService(enrollmentRepo, enrollmentCache, logger)
	testSvc := services.NewTestService(testRepo, testCache, logger)
	gradeSvc := services.NewGradeService(gradeRepo, gradeCache, logger)

	handler := http.NewHandler(userSvc, studentSvc, teacherSvc,
		courseSvc, enrollmentSvc, testSvc, gradeSvc, logger)

	addr := os.Getenv("SERVER_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8080"
	}
	logger.Info("trying to start server")
	if err := handler.Run(handler.Mount(), addr); err != nil {
		logger.Error("main: error while booting HTTP server", "error", err)
		os.Exit(1)
	}
	logger.Info("server has been shut down")
}
