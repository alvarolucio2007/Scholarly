package ports

import (
	"context"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
)

type StudentRepository interface {
	Repository[domain.Student]
	GetByEnrollment(ctx context.Context, enrollment string) (domain.Student, error)
}
