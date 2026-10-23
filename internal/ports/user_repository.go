package ports

import (
	"context"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
)

type UserFilter struct {
	Name  *string
	CPF   *string
	Email *string
}
type UserRepository interface {
	Repository[domain.User]
	List(ctx context.Context, filter UserFilter) ([]domain.User, error)
}
