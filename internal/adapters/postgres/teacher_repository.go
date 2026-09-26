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

var _ ports.TeacherRepository = (*TeacherRepo)(nil)

type TeacherRepo struct {
	db *sql.DB
}

func NewTeacherRepo(db *sql.DB) *TeacherRepo {
	return &TeacherRepo{db: db}
}

func (r *TeacherRepo) Create(ctx context.Context, teacher *domain.Teacher) error {
	query := `INSERT INTO teachers (user_id, department) VALUES ($1,$2)`

	ctx, cancel := context.WithTimeout(ctx, PostgresQueryTimeout)
	defer cancel()
	if _, err := r.db.ExecContext(ctx, query, teacher.UserID, teacher.Department); err != nil {
		translated := translateError(err)
		if !errors.Is(err, translated) {
			return translated
		}
		return fmt.Errorf("postgres: create teacher: %w", err)
	}
	return nil
}

func (r *TeacherRepo) GetByID(ctx context.Context, teacherID int64) (*domain.Teacher, error) {
	query := `SELECT user_id,department FROM teachers WHERE user_id=$1`
	var t domain.Teacher

	ctx, cancel := context.WithTimeout(ctx, PostgresQueryTimeout)
	defer cancel()
	err := r.db.QueryRowContext(ctx, query, teacherID).Scan(&t.UserID, &t.Department)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("postgres: getByID teacher: %w", err)
	}
	return &t, nil
}

func (r *TeacherRepo) List(ctx context.Context, filter ports.TeacherFilter) ([]*domain.Teacher, error) {
	wb := newWhereBuilder()
	if filter.Department != nil {
		wb.addILike("department", *filter.Department)
	}
	where, args := wb.build()
	query := `
    SELECT user_id,department
    FROM teachers
		` + where + `ORDER BY user_id ASC`
	ctx, cancel := context.WithTimeout(ctx, PostgresQueryTimeout)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("postgres: list teachers: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("postgres: rows.close list teachers: %v", err)
		}
	}()
	var teachers []*domain.Teacher
	for rows.Next() {
		var t domain.Teacher
		if err := rows.Scan(&t.UserID, &t.Department); err != nil {
			return nil, fmt.Errorf("postgres: scan teachers: %w", err)
		}
		teachers = append(teachers, &t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres: iterate teachers: %w", err)
	}
	return teachers, nil
}

func (r *TeacherRepo) Update(ctx context.Context, teacher *domain.Teacher) error {
	query := `UPDATE teachers
	SET
		department=COALESCE(NULLIF($1,''),department)
	WHERE user_id=$2`
	ctx, cancel := context.WithTimeout(ctx, PostgresQueryTimeout)
	defer cancel()
	res, err := r.db.ExecContext(ctx, query, teacher.Department, teacher.UserID)
	if err != nil {
		return fmt.Errorf("postgres: update teacher: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("postgres: rowsAffected update teacher: %w", err)
	}
	if count == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *TeacherRepo) Delete(ctx context.Context, teacherID int64) error {
	query := `DELETE FROM teachers WHERE user_id = $1`

	ctx, cancel := context.WithTimeout(ctx, PostgresQueryTimeout)
	defer cancel()

	res, err := r.db.ExecContext(ctx, query, teacherID)
	if err != nil {
		return fmt.Errorf("postgres: delete teacher: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("postgres: rowsAffected delete teacher: %w", err)
	}
	if count == 0 {
		return domain.ErrNotFound
	}
	return nil
}
