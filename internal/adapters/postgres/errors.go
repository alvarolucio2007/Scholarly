package postgres

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrNotFound            error = errors.New("not found")
	ErrUniqueViolation     error = errors.New("unique constraint violation")
	ErrForeignKeyViolation error = errors.New("foregin key violation")
)

func translateError(err error) error {
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			return ErrUniqueViolation
		case "23503": // foreign_key_violation
			return ErrForeignKeyViolation
		}
	}
	return err
}
