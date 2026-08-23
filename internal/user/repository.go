package user

import (
	"context"
)

type Repository interface {
	List(ctx context.Context) ([]User, error)
	FindByID(ctx context.Context, id string) (User, error)
}
