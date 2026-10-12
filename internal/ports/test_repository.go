package ports

import (
	"context"
	"time"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
)

type TestFilter struct {
	CourseID *int64
	Name     *string
	TestDate *time.Time
}
type TestRepository interface {
	Repository[domain.Test]
	List(ctx context.Context, filter TestFilter) ([]*domain.Test, error)
}
