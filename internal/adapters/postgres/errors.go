package postgres

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrUniqueViolation     error = errors.New("unique constraint violation")
	ErrForeignKeyViolation error = errors.New("foregin key violation")
)

func translateError(err error) error {
	if err == nil {
		return nil
	}

	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		switch pgErr.Code {
		case "23505": // unique_violation
			return ErrUniqueViolation
		case "23503": // foreign_key_violation
			return ErrForeignKeyViolation
		}
	}
	return err
}
