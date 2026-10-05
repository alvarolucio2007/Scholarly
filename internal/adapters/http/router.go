package http

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alvarolucio2007/Scholarly/docs"
	"github.com/alvarolucio2007/Scholarly/internal/env"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

const addr string = "localhost:8080"

func (h *Handler) Mount() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.ClientIPFromRemoteAddr)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{env.GetString("CORS_ALLOWED_ORIGIN", "http://localhost:5174")},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	r.Get("/health", h.healthCheckHandler)
	r.Get("/swagger/*", httpSwagger.Handler())
	// TODO: make it so that it gets user ID from path, all updates for that matter.(really gotta do this)
	r.Route("/users", func(r chi.Router) {
		r.Post("/", h.CreateUser)
		r.Get("/", h.ListUsers)
		r.Put("/", h.UpdateUser)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.GetUserByID)
			r.Delete("/", h.DeleteUser)
		})
	})
	r.Route("/students", func(r chi.Router) {
		r.Post("/", h.CreateStudent)
		r.Put("/", h.UpdateStudent)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.GetStudentByID)
			r.Delete("/", h.DeleteStudent)
		})
		r.Get("/enrollment/{id}", h.GetStudentByEnrollment)
	})
	r.Route("/teachers", func(r chi.Router) {
		r.Post("/", h.CreateTeacher)
		r.Put("/", h.UpdateTeacher)
		r.Get("/", h.ListTeachers)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.GetTeacherByID)
			r.Delete("/", h.DeleteTeacher)
		})
	})
	r.Route("/courses", func(r chi.Router) {
		r.Post("/", h.CreateCourse)
		r.Put("/", h.UpdateCourse)
		r.Get("/", h.ListCourses)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.GetCourseByID)
			r.Delete("/", h.DeleteCourse)
		})
		r.Route("/enroll/", func(r chi.Router) {
			r.Post("/", h.EnrollStudent)
		})
	})
	r.Route("/enrollments", func(r chi.Router) {
		r.Put("/", h.UpdateEnrollment)
		r.Get("/", h.ListEnrollments)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.GetEnrollmentByID)
			r.Delete("/", h.DeleteEnrollment)
		})
	})
	r.Route("/tests", func(r chi.Router) {
		r.Post("/", h.CreateTest)
		r.Put("/", h.UpdateTest)
		r.Get("/", h.ListTests)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.GetTestByID)
			r.Delete("/", h.DeleteTest)
		})
	})
	r.Route("/grades", func(r chi.Router) {
		r.Post("/", h.CreateGrade)
		r.Put("/", h.UpdateGrade)
		r.Get("/", h.ListGrades)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.GetGradeByID)
			r.Delete("/", h.DeleteGrade)
		})
		r.Route("/average", func(r chi.Router) {
			r.Get("/list", h.ListAverages)
			r.Get("/{student_id}/{course_id}", h.GetAverage)
		})
	})

	return r
}

func (h *Handler) Run(mux http.Handler) error {
	docs.SwaggerInfo.Version = "0.0.1"
	docs.SwaggerInfo.Host = addr
	srv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
	}
	shutdown := make(chan error, 1)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		shutdown <- srv.Shutdown(ctx)
	}()
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	err = <-shutdown
	if err != nil {
		return err
	}
	return nil
}
