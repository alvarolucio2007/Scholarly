package ports

import "context"

type Cache[T any] interface {
	Create(ctx context.Context, item *T) error
	Read(ctx context.Context, id int64) (*T, error)
	Delete(ctx context.Context, id int64) error
}
