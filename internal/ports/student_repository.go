package ports

import (
	"context"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
)

type StudentFilter struct {
	EnrollmentNumber *string
}
type StudentRepository interface {
	Repository[domain.Student]
	List(ctx context.Context, filter StudentFilter) ([]domain.Student, error)
}
