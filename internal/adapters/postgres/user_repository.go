package postgres

import (
	"database/sql"

	"github.com/alvarolucio2007/Scholarly/internal/ports"
)

var _ ports.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	db *sql.DB
}
