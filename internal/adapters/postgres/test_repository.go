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

var _ ports.TestRepository = (*TestRepo)(nil)

type TestRepo struct {
	db *sql.DB
}

func NewTestRepo(db *sql.DB) *TestRepo {
	return &TestRepo{db: db}
}

func (r *TestRepo) Create(ctx context.Context, test *domain.Test) error {
	query := `INSERT INTO tests (course_id,name,weight,test_date) VALUES ($1,$2,$3,$4) RETURNING id,created_at`
	ctx, cancel := context.WithTimeout(ctx, PostgresQueryTimeout)
	defer cancel()
	if err := r.db.QueryRowContext(ctx, query, test.CourseID, test.Name, test.Weight, test.TestDate).Scan(&test.ID, &test.CreatedAt); err != nil {
		translated := translateError(err)
		if !errors.Is(err, translated) {
			return translated
		}
		return fmt.Errorf("postgres: create test: %w", err)
	}
	return nil
}

func (r *TestRepo) GetByID(ctx context.Context, testID int64) (*domain.Test, error) {
	var t domain.Test
	query := `SELECT id,course_id,name,weight,test_date,created_at,updated_at FROM tests WHERE id=$1`
	ctx, cancel := context.WithTimeout(ctx, PostgresQueryTimeout)
	defer cancel()
	err := r.db.QueryRowContext(ctx, query, testID).Scan(&t.ID, &t.CourseID, &t.Name, &t.Weight, &t.TestDate, &t.CreatedAt, &t.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("postgres: get test by id: %w", err)
	}
	return &t, nil
}

func (r *TestRepo) List(ctx context.Context, filter ports.TestFilter) ([]*domain.Test, error) {
	wb := newWhereBuilder()
	if filter.CourseID != nil {
		wb.add("course_id", *filter.CourseID)
	}
	if filter.Name != nil {
		wb.addILike("name", *filter.Name)
	}
	if filter.TestDate != nil {
		wb.add("test_date", *filter.TestDate)
	}
	where, args := wb.build()
	query := `
    SELECT id,course_id,name,weight,test_date,created_at,updated_at
    FROM tests
` + where + ` ORDER BY created_at DESC`
	ctx, cancel := context.WithTimeout(ctx, PostgresQueryTimeout)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("postgres: list test: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("postgres: rows.close list test: %v", err)
		}
	}()
	var tests []*domain.Test
	for rows.Next() {
		var e domain.Test
		if err := rows.Scan(&e.ID, &e.CourseID, &e.Name, &e.Weight, &e.TestDate, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, fmt.Errorf("postgres: scan test: %w", err)
		}
		tests = append(tests, &e)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres: iterate test: %w", err)
	}
	return tests, nil
}

func (r *TestRepo) Update(ctx context.Context, test *domain.Test) error {
	query := `UPDATE tests
	SET
		course_id=COALESCE(NULLIF($1,0),course_id),
		name=COALESCE(NULLIF($2,''),name),
		weight=COALESCE(NULLIF($3,0),weight),
		test_date=COALESCE($4,test_date),
		updated_at = NOW()
	WHERE id=$5;
	`
	ctx, cancel := context.WithTimeout(ctx, PostgresQueryTimeout)
	defer cancel()
	res, err := r.db.ExecContext(ctx, query, test.CourseID, test.Name, test.Weight, test.TestDate, test.ID)
	if err != nil {
		return fmt.Errorf("postgres: update test: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("postgres: rows affected update test: %w", err)
	}
	if count == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *TestRepo) Delete(ctx context.Context, testID int64) error {
	query := `DELETE FROM tests WHERE id=$1`
	ctx, cancel := context.WithTimeout(ctx, PostgresQueryTimeout)
	defer cancel()
	res, err := r.db.ExecContext(ctx, query, testID)
	if err != nil {
		return fmt.Errorf("postgres: delete test: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("postgres: rows affected delete test: %w", err)
	}
	if count == 0 {
		return domain.ErrNotFound
	}
	return nil
}
