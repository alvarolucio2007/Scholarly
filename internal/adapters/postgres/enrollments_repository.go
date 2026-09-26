package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/alvarolucio2007/Scholarly/internal/ports"
)

var _ ports.EnrollmentRepository = (*EnrollmentRepo)(nil)

type EnrollmentRepo struct {
	db *sql.DB
}

func NewEnrollmentRepo(db *sql.DB) *EnrollmentRepo {
	return &EnrollmentRepo{db: db}
}

func (r *EnrollmentRepo) Create(ctx context.Context, enrollment *domain.Enrollment) error {
	query := `INSERT INTO enrollments (student_id,course_id,status) VALUES ($1,$2,$3) RETURNING id,enrolled_at`
	ctx, cancel := context.WithTimeout(ctx, PostgresQueryTimeout)
	defer cancel()
	if err := r.db.QueryRowContext(ctx, query, enrollment.StudentID, enrollment.CourseID, enrollment.Status).Scan(&enrollment.ID, &enrollment.EnrolledAt); err != nil {
		translated := translateError(err)
		if !errors.Is(err, translated) {
			return translated
		}
		return fmt.Errorf("postgres: create enrollment: %w", err)
	}
	return nil
}

func (r *EnrollmentRepo) GetByID(ctx context.Context, enrollmentID int64) (*domain.Enrollment, error) {
	var e domain.Enrollment
	query := `SELECT id,student_id,course_id,enrolled_at,status FROM enrollments WHERE id=$1`
	ctx, cancel := context.WithTimeout(ctx, PostgresQueryTimeout)
	defer cancel()
	err := r.db.QueryRowContext(ctx, query, enrollmentID).Scan(&e.ID, &e.StudentID, &e.CourseID, &e.EnrolledAt, &e.Status)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("postgres: get enrollment by id: %w", err)
	}
	return &e, nil
}

func (r *EnrollmentRepo) List(ctx context.Context, filter ports.EnrollmentFilter) ([]*domain.Enrollment, error) {
	wb := newWhereBuilder()
	if filter.StudentID != nil {
		wb.add("student_id", *filter.StudentID)
	}
	if filter.CourseID != nil {
		wb.add("course_id", *filter.CourseID)
	}
	if filter.Status != nil {
		wb.add("status", *filter.Status)
	}
	where, args := wb.build()
	query := `
    SELECT id,student_id,course_id,enrolled_at,status
    FROM enrollments
` + where + `ORDER BY enrolled_at DESC`
	ctx, cancel := context.WithTimeout(ctx, PostgresQueryTimeout)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("postgres: list enrollment: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("postgres: rows.close list enrollment: %v", err)
		}
	}()
	var enrollments []*domain.Enrollment
	for rows.Next() {
		var e domain.Enrollment
		if err := rows.Scan(&e.ID, &e.StudentID, &e.CourseID, &e.EnrolledAt, &e.Status); err != nil {
			return nil, fmt.Errorf("postgres: scan enrollments: %w", err)
		}
		enrollments = append(enrollments, &e)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres: iterate enrollments: %w", err)
	}
	return enrollments, nil
}

func (r *EnrollmentRepo) Update(ctx context.Context, enrollment *domain.Enrollment) error {
	query := `UPDATE enrollments
	SET
		student_id=COALESCE(NULLIF($1,0),student_id),
		course_id=COALESCE(NULLIF($2,0),course_id)
	WHERE id=$3;
	`
	ctx, cancel := context.WithTimeout(ctx, PostgresQueryTimeout)
	defer cancel()
	res, err := r.db.ExecContext(ctx, query, enrollment.StudentID, enrollment.CourseID, enrollment.ID)
	if err != nil {
		return fmt.Errorf("postgres: update enrollment: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("postgres: rows affected update enrollment: %w", err)
	}
	if count == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *EnrollmentRepo) Delete(ctx context.Context, enrollmentID int64) error {
	query := `DELETE FROM enrollments WHERE id=$1`
	ctx, cancel := context.WithTimeout(ctx, PostgresQueryTimeout)
	defer cancel()
	res, err := r.db.ExecContext(ctx, query, enrollmentID)
	if err != nil {
		return fmt.Errorf("postgres: delete enrollment: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("postgres: rows affected delete enrollments: %w", err)
	}
	if count == 0 {
		return domain.ErrNotFound
	}
	return nil
}
