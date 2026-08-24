package project

import (
	"context"
	"strings"
)

type Service struct {
	repository Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repository: repo,
	}
}

func (s *Service) CreateProject(ctx context.Context, input CreateProjectInput) (*Project, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, ErrNameMustPresent
	}

	createdBy := strings.TrimSpace(input.CreatedBy)
	if createdBy == "" {
		return nil, ErrCreatedByMustPresent
	}

	project := Project{
		Name:        name,
		Description: strings.TrimSpace(input.Description),
		CreatedBy:   createdBy,
	}

	createdProject, err := s.repository.Create(ctx, project)
	if err != nil {
		return nil, err
	}

	return createdProject, nil
}

func (s *Service) GetProjects(ctx context.Context) ([]Project, error) {
	return s.repository.List(ctx)
}

func (s *Service) GetProjectByID(ctx context.Context, id string) (*Project, error) {
	return s.repository.FindByID(ctx, id)
}
