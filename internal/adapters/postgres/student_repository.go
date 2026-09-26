package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/alvarolucio2007/Scholarly/internal/ports"
)

var _ ports.StudentRepository = (*StudentRepo)(nil)

type StudentRepo struct {
	db *sql.DB
}

func NewStudentRepo(db *sql.DB) *StudentRepo {
	return &StudentRepo{db: db}
}

func (r *StudentRepo) Create(ctx context.Context, student *domain.Student) error {
	query := `INSERT INTO students (user_id,enrollment_number) VALUES ($1,$2) RETURNING user_id`
	ctx, cancel := context.WithTimeout(ctx, PostgresQueryTimeout)
	defer cancel()
	if err := r.db.QueryRowContext(ctx, query, student.UserID, student.EnrollmentNumber).Scan(&student.UserID); err != nil {
		translated := translateError(err)
		if translated != err {
			return translated
		}
		return fmt.Errorf("postgres: create student: %w", err)
	}
	return nil
}

func (r *StudentRepo) GetByID(ctx context.Context, studentID int64) (*domain.Student, error) {
	query := `SELECT user_id,enrollment_number FROM students WHERE user_id=$1`
	var s domain.Student

	ctx, cancel := context.WithTimeout(ctx, PostgresQueryTimeout)
	defer cancel()
	err := r.db.QueryRowContext(ctx, query, studentID).Scan(&s.UserID, &s.EnrollmentNumber)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("postgres: getByID student: %w", err)
	}
	return &s, nil
}

func (r *StudentRepo) GetByEnrollment(ctx context.Context, enrollment string) (*domain.Student, error) {
	query := `SELECT user_id,enrollment_number FROM students WHERE enrollment_number=$1`
	var s domain.Student

	ctx, cancel := context.WithTimeout(ctx, PostgresQueryTimeout)
	defer cancel()
	err := r.db.QueryRowContext(ctx, query, enrollment).Scan(&s.UserID, &s.EnrollmentNumber)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("postgres: getByEnrollment student: %w", err)
	}
	return &s, nil
}

func (r *StudentRepo) Update(ctx context.Context, student *domain.Student) error {
	query := `UPDATE students
	SET
		enrollment_number=COALESCE(NULLIF($1,''),enrollment_number),
	WHERE user_id=$2`
	ctx, cancel := context.WithTimeout(ctx, PostgresQueryTimeout)
	defer cancel()
	res, err := r.db.ExecContext(ctx, query, student.EnrollmentNumber, student.UserID)
	if err != nil {
		return fmt.Errorf("postgres: update student: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("postgres: rowsAffected update student: %w", err)
	}
	if count == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *StudentRepo) Delete(ctx context.Context, studentID int64) error {
	query := `DELETE FROM students WHERE user_id = $1`

	ctx, cancel := context.WithTimeout(ctx, PostgresQueryTimeout)
	defer cancel()

	res, err := r.db.ExecContext(ctx, query, studentID)
	if err != nil {
		return fmt.Errorf("postgres: delete student: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("postgres: rowsAffected delete student: %w", err)
	}
	if count == 0 {
		return domain.ErrNotFound
	}
	return nil
}
