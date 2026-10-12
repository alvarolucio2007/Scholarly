package ports

import (
	"context"
)

type Repository[T any] interface {
	Create(ctx context.Context, item *T) error
	GetByID(ctx context.Context, id int64) (*T, error)
	Update(ctx context.Context, item *T) error
	Delete(ctx context.Context, id int64) error
}
