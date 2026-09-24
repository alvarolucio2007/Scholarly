package domain

import (
	"errors"
)

var ErrDepartmentRequired = errors.New("department is required")

type Teacher struct {
	UserID     int64
	Department string
}
