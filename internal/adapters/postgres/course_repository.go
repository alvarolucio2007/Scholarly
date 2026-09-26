package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/alvarolucio2007/Scholarly/internal/ports"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var _ ports.CourseRepository = (*CourseRepo)(nil)

type CourseRepo struct {
	db *sql.DB
}

func NewCourseRepo(db *sql.DB) *CourseRepo {
	return &CourseRepo{db: db}
}

func (r *CourseRepo) Create(ctx context.Context, course *domain.Course) error {
	query := `INSERT INTO courses (teacher_id,name,code,semester,max_students) VALUES ($1,$2,$3,$4,$5) RETURNING id,created_at`
	ctx, cancel := context.WithTimeout(ctx, PostgresQueryTimeout)
	defer cancel()
	if err := r.db.QueryRowContext(ctx, query, course.TeacherID, course.Name, course.Code, course.Semester, course.MaxStudents).Scan(&course.ID, &course.CreatedAt); err != nil {
		translated := translateError(err)
		if !errors.Is(err, translated) {
			return translated
		}
		return fmt.Errorf("postgres: create course: %w", err)
	}
	return nil
}

func (r *CourseRepo) GetByID(ctx context.Context, courseID int64) (*domain.Course, error) {
	var c domain.Course
	query := `SELECT id,teacher_id,name,code,semester,max_students,created_at,updated_at FROM courses WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, PostgresQueryTimeout)
	defer cancel()
	err := r.db.QueryRowContext(ctx, query, courseID).Scan(&c.ID, &c.TeacherID, &c.Name, &c.Code, &c.Semester, &c.MaxStudents, &c.CreatedAt, &c.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("postgres: get course by id: %w", err)
	}
	return &c, nil
}

func (r *CourseRepo) List(ctx context.Context, filter ports.CourseFilter) ([]*domain.Course, error) {
	wb := newWhereBuilder()
	if filter.TeacherID != nil {
		wb.add("teacher_id", *filter.TeacherID)
	}
	if filter.Name != nil {
		wb.addILike("name", *filter.Name)
	}
	if filter.Code != nil {
		wb.addILike("code", *filter.Code)
	}
	if filter.Semester != nil {
		wb.addILike("semester", *filter.Semester)
	}
	where, args := wb.build()
	query := `
    SELECT id,teacher_id,name,code,semester,max_students,created_at,updated_at
    FROM courses
` + where + `
    ORDER BY created_at DESC
`
	ctx, cancel := context.WithTimeout(ctx, PostgresQueryTimeout)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("postgres: list courses: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("postgres: rows.close list courses: %v", err)
		}
	}()
	var courses []*domain.Course
	for rows.Next() {
		var c domain.Course
		if err := rows.Scan(&c.ID, &c.TeacherID, &c.Name, &c.Code, &c.Semester, &c.MaxStudents, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("postgres: scan courses: %w", err)
		}
		courses = append(courses, &c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres: iterate courses: %w", err)
	}
	return courses, nil
}

func (r *CourseRepo) Update(ctx context.Context, course *domain.Course) error {
	query := ` UPDATE courses
	SET
		teacher_id=COALESCE(NULLIF($1,0),teacher_id),
		name=COALESCE(NULLIF($2,''),name),
		code=COALESCE(NULLIF($3,''),code),
		semester=COALESCE(NULLIF($4,''),semester),
		max_students=COALESCE(NULLIF($5,0),max_students),
		updated_at=NOW()
	WHERE id=$6;
	`
	ctx, cancel := context.WithTimeout(ctx, PostgresQueryTimeout)
	defer cancel()
	res, err := r.db.ExecContext(ctx, query, course.TeacherID, course.Name, course.Code, course.Semester, course.MaxStudents, course.ID)
	if err != nil {
		translated := translateError(err)
		if !errors.Is(err, translated) {
			return translated
		}
		return fmt.Errorf("postgres: update courses: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("postgres: rows affected update courses: %w", err)
	}
	if count == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *CourseRepo) Delete(ctx context.Context, courseID int64) error {
	query := `DELETE FROM courses WHERE id=$1`
	ctx, cancel := context.WithTimeout(ctx, PostgresQueryTimeout)
	defer cancel()
	res, err := r.db.ExecContext(ctx, query, courseID)
	if err != nil {
		return fmt.Errorf("postgres: delete course: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("postgres: rows affected delete course: %w", err)
	}
	if count == 0 {
		return domain.ErrNotFound
	}
	return nil
}
