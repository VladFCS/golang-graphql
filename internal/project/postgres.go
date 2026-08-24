package project

import (
	"context"
	"database/sql"
	"errors"

	projectdb "github.com/vladfc/ghira/internal/database/db"
)

type PostgresRepository struct {
	queries *projectdb.Queries
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		queries: projectdb.New(db),
	}
}

func (r *PostgresRepository) Create(ctx context.Context, project Project) (*Project, error) {
	row, err := r.queries.CreateProject(ctx, toCreateProjectParams(project))
	if err != nil {
		return nil, err
	}

	createdProject := toProject(row)
	return &createdProject, nil
}

func (r *PostgresRepository) List(ctx context.Context) ([]Project, error) {
	rows, err := r.queries.ListProjects(ctx)
	if err != nil {
		return nil, err
	}

	return toProjects(rows), nil
}

func (r *PostgresRepository) FindByID(ctx context.Context, id string) (*Project, error) {
	row, err := r.queries.GetProjectByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	foundProject := toProject(row)
	return &foundProject, nil
}

func toCreateProjectParams(project Project) projectdb.CreateProjectParams {
	return projectdb.CreateProjectParams{
		Name:        project.Name,
		Description: project.Description,
		CreatedBy:   project.CreatedBy,
	}
}

func toProject(row projectdb.Project) Project {
	return Project{
		ID:          row.ID,
		Name:        row.Name,
		Description: row.Description,
		CreatedBy:   row.CreatedBy,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}

func toProjects(rows []projectdb.Project) []Project {
	projects := make([]Project, 0, len(rows))

	for _, row := range rows {
		projects = append(projects, toProject(row))
	}

	return projects
}

var _ Repository = (*PostgresRepository)(nil)
