package project

import "context"

type Repository interface {
	Create(ctx context.Context, project Project) (*Project, error)
	List(ctx context.Context) ([]Project, error)
	FindByID(ctx context.Context, id string) (*Project, error)
}

