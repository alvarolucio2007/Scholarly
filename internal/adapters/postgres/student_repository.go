package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/alvarolucio2007/Scholarly/internal/ports"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	_ ports.StudentRepository          = (*StudentRepo)(nil)
	_ ports.Repository[domain.Student] = (*StudentRepo)(nil)
)

type StudentRepo struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewStudentRepo(db *sql.DB) *StudentRepo {
	return &StudentRepo{db: db}
}

func (r *StudentRepo) Create(ctx context.Context, student *domain.Student) error {
	query := `INSERT INTO students (user_id,enrollment_number) VALUES ($1,$2)`
	ctx, cancel := context.WithTimeout(ctx, PostgresQueryTimeout)
	defer cancel()
	if _, err := r.db.ExecContext(ctx, query, student.UserID, student.EnrollmentNumber); err != nil {
		translated := translateError(err)
		if !errors.Is(err, translated) {
			if errors.Is(translated, domain.ErrForeignKeyViolation) {
				return domain.ErrUserNotFound
			}
			return translated
		}
		return fmt.Errorf("postgres: create student: %w", err)
	}
	return nil
}

func (r *StudentRepo) GetByID(ctx context.Context, studentID int64) (domain.Student, error) {
	query := `SELECT user_id,enrollment_number FROM students WHERE user_id=$1`
	var s domain.Student

	ctx, cancel := context.WithTimeout(ctx, PostgresQueryTimeout)
	defer cancel()
	err := r.db.QueryRowContext(ctx, query, studentID).Scan(&s.UserID, &s.EnrollmentNumber)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Student{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Student{}, fmt.Errorf("postgres: getByID student: %w", err)
	}
	return s, nil
}

// List implements [ports.StudentRepository].
func (r *StudentRepo) List(ctx context.Context, filter ports.StudentFilter) ([]domain.Student, error) {
	wb := newWhereBuilder()
	if filter.EnrollmentNumber != nil {
		wb.addILike("enrollment_number", *filter.EnrollmentNumber)
	}
	where, args := wb.build()
	query := `
    SELECT user_id,enrollment_number
    FROM students
		` + where + `ORDER BY user_id ASC`
	ctx, cancel := context.WithTimeout(ctx, PostgresQueryTimeout)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("postgres: list students: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			r.logger.ErrorContext(ctx, "postgres: rows.close list students ", "error", err)
		}
	}()
	var teachers []domain.Student
	for rows.Next() {
		var t domain.Student
		if err := rows.Scan(&t.UserID, &t.EnrollmentNumber); err != nil {
			return nil, fmt.Errorf("postgres: scan students: %w", err)
		}
		teachers = append(teachers, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres: iterate teachers: %w", err)
	}
	return teachers, nil
}

func (r *StudentRepo) Update(ctx context.Context, student domain.Student) error {
	query := `UPDATE students
	SET
		enrollment_number=COALESCE(NULLIF($1,''),enrollment_number)
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
