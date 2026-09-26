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

var _ ports.GradeRepository = (*GradeRepo)(nil)

type GradeRepo struct {
	db *sql.DB
}

func (r *GradeRepo) Create(ctx context.Context, grade *domain.Grade) error {
	query := `INSERT INTO grades (test_id,enrollment_id,value) VALUES ($1,$2,$3) RETURNING id,created_at`
	ctx, cancel := context.WithTimeout(ctx, PostgresQueryTimeout)
	defer cancel()
	if err := r.db.QueryRowContext(ctx, query, grade.TestID, grade.EnrollmentID, grade.Value).Scan(&grade.ID, &grade.CreatedAt); err != nil {
		translated := translateError(err)
		if !errors.Is(err, translated) {
			return translated
		}
		return fmt.Errorf("postgres: create grade: %w", err)
	}
	return nil
}

func (r *GradeRepo) GetByID(ctx context.Context, gradeID int64) (*domain.Grade, error) {
	var g domain.Grade
	query := `SELECT id,test_id,enrollment_id,value,created_at,updated_at FROM grades WHERE id=$1`
	ctx, cancel := context.WithTimeout(ctx, PostgresQueryTimeout)
	defer cancel()
	err := r.db.QueryRowContext(ctx, query, gradeID).Scan(&g.ID, &g.TestID, &g.EnrollmentID, &g.Value, &g.CreatedAt, &g.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("postgres: get grade by id: %w", err)
	}
	return &g, nil
}

func (r *GradeRepo) List(ctx context.Context, filter ports.GradeFilter) ([]*domain.Grade, error) {
	wb := newWhereBuilder()
	if filter.TestID != nil {
		wb.add("test_id", *filter.TestID)
	}
	if filter.EnrollmentID != nil {
		wb.add("enrollment_id", *filter.EnrollmentID)
	}
	where, args := wb.build()
	query := `
		SELECT id,test_id,enrollment_id,value,created_at,updated_at
		FROM grades
` + where + ` ORDER BY created_at DESC`
	ctx, cancel := context.WithTimeout(ctx, PostgresQueryTimeout)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("postgres: list grade: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("postgres: rows.close list grade: %v", err)
		}
	}()
	var grades []*domain.Grade
	for rows.Next() {
		var g domain.Grade
		if err := rows.Scan(&g.ID, &g.TestID, &g.EnrollmentID, &g.Value, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return nil, fmt.Errorf("postgres: scan grade: %w", err)
		}
		grades = append(grades, &g)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres: iterate grade: %w", err)
	}
	return grades, nil
}

func (r *GradeRepo) Update(ctx context.Context, grade *domain.Grade) error {
	query := `UPDATE grades
	SET
		test_id=COALESCE(NULLIF($1,0),test_id),
		enrollment_id=COALESCE(NULLIF($2,0),enrollment_id),
		value=COALESCE($3,value),
		updated_at = NOW()
	WHERE id=$4;
	`
	ctx, cancel := context.WithTimeout(ctx, PostgresQueryTimeout)
	defer cancel()
	res, err := r.db.ExecContext(ctx, query, grade.TestID, grade.EnrollmentID, grade.Value, grade.ID)
	if err != nil {
		return fmt.Errorf("postgres: update grade: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("postgres: rows affected update grade: %w", err)
	}
	if count == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *GradeRepo) Delete(ctx context.Context, gradeID int64) error {
	query := `DELETE FROM grades WHERE id=$1`
	ctx, cancel := context.WithTimeout(ctx, PostgresQueryTimeout)
	defer cancel()
	res, err := r.db.ExecContext(ctx, query, gradeID)
	if err != nil {
		return fmt.Errorf("postgres: delete grade: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("postgres: rows affected delete grade: %w", err)
	}
	if count == 0 {
		return domain.ErrNotFound
	}
	return nil
}
