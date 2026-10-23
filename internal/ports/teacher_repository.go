package ports

import (
	"context"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
)

type TeacherFilter struct {
	Department *string
}
type TeacherRepository interface {
	Repository[domain.Teacher]
	List(ctx context.Context, filter TeacherFilter) ([]domain.Teacher, error)
}
