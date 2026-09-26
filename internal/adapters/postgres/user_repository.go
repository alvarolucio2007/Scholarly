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

var _ ports.UserRepository = (*UserRepo)(nil)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (name,cpf,email,password_hash) VALUES ($1,$2,$3,$4) RETURNING id,created_at`
	ctx, cancel := context.WithTimeout(ctx, PostgresQueryTimeout)
	defer cancel()
	if err := r.db.QueryRowContext(ctx, query, user.Name, user.CPF, user.Email, user.PasswordHash).Scan(&user.ID, &user.CreatedAt); err != nil {
		return fmt.Errorf("postgres: create user: %w", err)
	}
	return nil
}

func (r *UserRepo) GetByID(ctx context.Context, userID int64) (*domain.User, error) {
	var u domain.User
	query := `SELECT id,name,cpf,email,created_at,updated_at FROM users WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, PostgresQueryTimeout)
	defer cancel()
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&u.ID, &u.Name, &u.CPF, &u.Email, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("postgres: get user by id %w", err)
	}
	return &u, nil
}

func (r *UserRepo) List(ctx context.Context, filter ports.UserFilter) ([]*domain.User, error) {
	wb := newWhereBuilder()
	if filter.Name != nil {
		wb.addILike("name", *filter.Name)
	}
	if filter.CPF != nil {
		wb.add("cpf", *filter.CPF)
	}
	if filter.Email != nil {
		wb.addILike("email", *filter.Email)
	}
	where, args := wb.build()
	query := `
    SELECT id, name, cpf, email, created_at, updated_at
    FROM users
` + where + `
    ORDER BY created_at DESC
`
	ctx, cancel := context.WithTimeout(ctx, PostgresQueryTimeout)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("postgres: list users: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("postgres: rows.close list users: %v", err)
		}
	}()
	var users []*domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Name, &u.CPF, &u.Email, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, fmt.Errorf("postgres: scan users: %w", err)
		}
		users = append(users, &u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres: iterate users: %w", err)
	}
	return users, nil
}

func (r *UserRepo) Update(ctx context.Context, user *domain.User) error {
	query := `UPDATE users
	SET
		name=COALESCE(NULLIF($1,''),name),
		cpf=COALESCE(NULLIF($2,''),cpf),
		email=COALESCE(NULLIF($3,''),email),
		password_hash=COALESCE($4,password_hash)
	WHERE id=$5;
	`
	ctx, cancel := context.WithTimeout(ctx, PostgresQueryTimeout)
	defer cancel()
	res, err := r.db.ExecContext(ctx, query, user.Name, user.CPF, user.Email, user.PasswordHash, user.ID)
	if err != nil {
		return fmt.Errorf("postgres: update users: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("postgres: rows affected update users: %w", err)
	}
	if count == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *UserRepo) Delete(ctx context.Context, userID int64) error {
	query := `DELETE FROM users WHERE id=$1`
	ctx, cancel := context.WithTimeout(ctx, PostgresQueryTimeout)
	defer cancel()
	res, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("postgres: delete user: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("postgres: rows affected delete users: %w", err)
	}
	if count == 0 {
		return domain.ErrNotFound
	}
	return nil
}
