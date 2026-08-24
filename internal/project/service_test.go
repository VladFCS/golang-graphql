package project

import (
	"context"
	"errors"
	"testing"
	"time"
)

type mockRepository struct {
	createCalled  bool
	createProject Project
	createResult  *Project
	createErr     error

	listCalled bool
	listResult []Project
	listErr    error

	findCalled bool
	findID     string
	findResult *Project
	findErr    error
}

func (r *mockRepository) Create(ctx context.Context, project Project) (*Project, error) {
	r.createCalled = true
	r.createProject = project
	return r.createResult, r.createErr
}

func (r *mockRepository) List(ctx context.Context) ([]Project, error) {
	r.listCalled = true
	return r.listResult, r.listErr
}

func (r *mockRepository) FindByID(ctx context.Context, id string) (*Project, error) {
	r.findCalled = true
	r.findID = id
	return r.findResult, r.findErr
}

func TestServiceCreateProject(t *testing.T) {
	t.Parallel()

	expectedProject := newTestProject(
		"f4fc94ea-2e2c-44a4-9ed0-e16e404ef873",
		"Issue Tracker",
		"GraphQL backend learning project",
		"8b7d8a0a-6c7e-4bb2-bb6a-68c73a4b6b01",
	)
	repo := &mockRepository{
		createResult: &expectedProject,
	}
	service := NewService(repo)

	project, err := service.CreateProject(context.Background(), CreateProjectInput{
		Name:        " Issue Tracker ",
		Description: " GraphQL backend learning project ",
		CreatedBy:   " 8b7d8a0a-6c7e-4bb2-bb6a-68c73a4b6b01 ",
	})
	if err != nil {
		t.Fatalf("CreateProject returned error: %v", err)
	}

	if !repo.createCalled {
		t.Fatal("expected repository Create to be called")
	}

	if repo.createProject.Name != "Issue Tracker" {
		t.Fatalf("expected project name to be trimmed, got %q", repo.createProject.Name)
	}

	if repo.createProject.Description != "GraphQL backend learning project" {
		t.Fatalf("expected project description to be trimmed, got %q", repo.createProject.Description)
	}

	if repo.createProject.CreatedBy != "8b7d8a0a-6c7e-4bb2-bb6a-68c73a4b6b01" {
		t.Fatalf("expected project creator to be trimmed, got %q", repo.createProject.CreatedBy)
	}

	if project == nil {
		t.Fatal("expected created project, got nil")
	}

	if *project != expectedProject {
		t.Fatalf("expected created project %+v, got %+v", expectedProject, *project)
	}
}

func TestServiceCreateProjectRequiresName(t *testing.T) {
	t.Parallel()

	repo := &mockRepository{}
	service := NewService(repo)

	project, err := service.CreateProject(context.Background(), CreateProjectInput{
		Name:      " ",
		CreatedBy: "8b7d8a0a-6c7e-4bb2-bb6a-68c73a4b6b01",
	})
	if !errors.Is(err, ErrNameMustPresent) {
		t.Fatalf("expected error %v, got %v", ErrNameMustPresent, err)
	}

	if project != nil {
		t.Fatalf("expected nil project, got %+v", project)
	}

	if repo.createCalled {
		t.Fatal("expected repository Create not to be called")
	}
}

func TestServiceCreateProjectRequiresCreatedBy(t *testing.T) {
	t.Parallel()

	repo := &mockRepository{}
	service := NewService(repo)

	project, err := service.CreateProject(context.Background(), CreateProjectInput{
		Name:      "Issue Tracker",
		CreatedBy: " ",
	})
	if !errors.Is(err, ErrCreatedByMustPresent) {
		t.Fatalf("expected error %v, got %v", ErrCreatedByMustPresent, err)
	}

	if project != nil {
		t.Fatalf("expected nil project, got %+v", project)
	}

	if repo.createCalled {
		t.Fatal("expected repository Create not to be called")
	}
}

func TestServiceCreateProjectReturnsRepositoryError(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("create project failed")
	repo := &mockRepository{
		createErr: expectedErr,
	}
	service := NewService(repo)

	project, err := service.CreateProject(context.Background(), CreateProjectInput{
		Name:      "Issue Tracker",
		CreatedBy: "8b7d8a0a-6c7e-4bb2-bb6a-68c73a4b6b01",
	})
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}

	if project != nil {
		t.Fatalf("expected nil project, got %+v", project)
	}
}

func TestServiceGetProjects(t *testing.T) {
	t.Parallel()

	expectedProjects := []Project{
		newTestProject("f4fc94ea-2e2c-44a4-9ed0-e16e404ef873", "Issue Tracker", "", "8b7d8a0a-6c7e-4bb2-bb6a-68c73a4b6b01"),
	}
	repo := &mockRepository{
		listResult: expectedProjects,
	}
	service := NewService(repo)

	projects, err := service.GetProjects(context.Background())
	if err != nil {
		t.Fatalf("GetProjects returned error: %v", err)
	}

	if !repo.listCalled {
		t.Fatal("expected repository List to be called")
	}

	if len(projects) != len(expectedProjects) {
		t.Fatalf("expected %d projects, got %d", len(expectedProjects), len(projects))
	}
}

func TestServiceGetProjectByID(t *testing.T) {
	t.Parallel()

	expectedProject := newTestProject("f4fc94ea-2e2c-44a4-9ed0-e16e404ef873", "Issue Tracker", "", "8b7d8a0a-6c7e-4bb2-bb6a-68c73a4b6b01")
	repo := &mockRepository{
		findResult: &expectedProject,
	}
	service := NewService(repo)

	project, err := service.GetProjectByID(context.Background(), expectedProject.ID)
	if err != nil {
		t.Fatalf("GetProjectByID returned error: %v", err)
	}

	if !repo.findCalled {
		t.Fatal("expected repository FindByID to be called")
	}

	if repo.findID != expectedProject.ID {
		t.Fatalf("expected repository FindByID id %q, got %q", expectedProject.ID, repo.findID)
	}

	if project == nil {
		t.Fatal("expected project, got nil")
	}
}

func newTestProject(id, name, description, createdBy string) Project {
	now := time.Date(2026, 8, 24, 10, 0, 0, 0, time.UTC)

	return Project{
		ID:          id,
		Name:        name,
		Description: description,
		CreatedBy:   createdBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
