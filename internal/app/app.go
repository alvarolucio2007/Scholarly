package app

import (
	"database/sql"
	"log/slog"

	"github.com/alvarolucio2007/Scholarly/internal/adapters/argon2"
	"github.com/alvarolucio2007/Scholarly/internal/adapters/http"
	"github.com/alvarolucio2007/Scholarly/internal/adapters/postgres"
	"github.com/alvarolucio2007/Scholarly/internal/adapters/valkey"
	"github.com/alvarolucio2007/Scholarly/internal/services"
	"github.com/redis/go-redis/v9"
)

type App struct {
	Handler *http.Handler
}

func New(db *sql.DB, rdb *redis.Client, logger *slog.Logger) *App {
	userRepo := postgres.NewUserRepo(db, logger)
	hasher := argon2.NewArgon2Hasher(64*1024, 3, 2)
	userCache := valkey.NewUserCache(rdb)
	userSvc := services.NewUserService(userRepo, hasher, userCache, logger)

	studentRepo := postgres.NewStudentRepo(db)
	studentCache := valkey.NewStudentCache(rdb)
	studentSvc := services.NewStudentService(studentRepo, studentCache, logger)

	teacherRepo := postgres.NewTeacherRepo(db, logger)
	teacherCache := valkey.NewTeacherCache(rdb)
	teacherSvc := services.NewTeacherService(teacherRepo, teacherCache, logger)

	courseRepo := postgres.NewCourseRepo(db, logger)
	courseCache := valkey.NewCourseCache(rdb)
	courseSvc := services.NewCourseService(courseRepo, courseCache, logger)

	enrollmentRepo := postgres.NewEnrollmentRepo(db, logger)
	enrollmentCache := valkey.NewEnrollmentCache(rdb)
	enrollmentSvc := services.NewEnrollmentService(enrollmentRepo, enrollmentCache, logger)

	testRepo := postgres.NewTestRepo(db, logger)
	testCache := valkey.NewTestCache(rdb)
	testSvc := services.NewTestService(testRepo, testCache, logger)

	gradeRepo := postgres.NewGradeRepo(db, logger)
	gradeCache := valkey.NewGradeCache(rdb)
	gradeSvc := services.NewGradeService(gradeRepo, gradeCache, logger)

	handler := http.NewHandler(userSvc, studentSvc, teacherSvc, courseSvc, enrollmentSvc, testSvc, gradeSvc, logger)
	return &App{Handler: handler}
}
