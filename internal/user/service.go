package user

import "context"

type Service struct {
	repository Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repository: repo,
	}
}

func (s *Service) GetUsers(ctx context.Context) ([]User, error) {
	return s.repository.List(ctx)
}

func (s *Service) GetUserByID(ctx context.Context, id string) (User, error) {
	return s.repository.FindByID(ctx, id)
}
