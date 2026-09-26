package postgres

import (
	"errors"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
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
			switch pgErr.ConstraintName {
			case "users_email_key":
				return domain.ErrEmailAlreadyExists
			case "users_cpf_key":
				return domain.ErrCPFAlreadyExists
			case "students_enrollment_number_key":
				return domain.ErrEnrollmentAlreadyExists
			case "courses_code_semester_unique":
				return domain.ErrCourseCodeAlreadyExists
			}
			return domain.ErrConflict
		case "23503":
			return domain.ErrForeignKeyViolation
		}
	}
	return err
}
