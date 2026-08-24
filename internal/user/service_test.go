package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type mockRepository struct {
	createCalled bool
	createUser   User
	createResult *User
	createErr    error

	listCalled bool
	listResult []User
	listErr    error

	findCalled bool
	findID     string
	findResult *User
	findErr    error
}

func (r *mockRepository) Create(ctx context.Context, user User) (*User, error) {
	r.createCalled = true
	r.createUser = user
	return r.createResult, r.createErr
}

func (r *mockRepository) List(ctx context.Context) ([]User, error) {
	r.listCalled = true
	return r.listResult, r.listErr
}

func (r *mockRepository) FindByID(ctx context.Context, id string) (*User, error) {
	r.findCalled = true
	r.findID = id
	return r.findResult, r.findErr
}

func TestServiceGetUsers(t *testing.T) {
	t.Parallel()

	expectedUsers := []User{
		newTestUser("8b7d8a0a-6c7e-4bb2-bb6a-68c73a4b6b01", "olena.koval@example.com", "olena"),
		newTestUser("1c4a3df7-4a5d-41c5-8708-3a8a94f52c21", "maksym.shevchenko@example.com", "maksym"),
	}
	repo := &mockRepository{
		listResult: expectedUsers,
	}
	service := NewService(repo)

	users, err := service.GetUsers(context.Background())
	if err != nil {
		t.Fatalf("GetUsers returned error: %v", err)
	}

	if !repo.listCalled {
		t.Fatal("expected repository List to be called")
	}

	if len(users) != len(expectedUsers) {
		t.Fatalf("expected %d users, got %d", len(expectedUsers), len(users))
	}

	for i, expectedUser := range expectedUsers {
		if users[i] != expectedUser {
			t.Fatalf("expected user at index %d to be %+v, got %+v", i, expectedUser, users[i])
		}
	}
}

func TestServiceGetUsersReturnsRepositoryError(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("list users failed")
	repo := &mockRepository{
		listErr: expectedErr,
	}
	service := NewService(repo)

	users, err := service.GetUsers(context.Background())
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}

	if users != nil {
		t.Fatalf("expected nil users on error, got %+v", users)
	}
}

func TestServiceGetUserByID(t *testing.T) {
	t.Parallel()

	expectedUser := newTestUser("8b7d8a0a-6c7e-4bb2-bb6a-68c73a4b6b01", "olena.koval@example.com", "olena")
	repo := &mockRepository{
		findResult: &expectedUser,
	}
	service := NewService(repo)

	user, err := service.GetUserByID(context.Background(), expectedUser.ID)
	if err != nil {
		t.Fatalf("GetUserByID returned error: %v", err)
	}

	if !repo.findCalled {
		t.Fatal("expected repository FindByID to be called")
	}

	if repo.findID != expectedUser.ID {
		t.Fatalf("expected repository FindByID id %q, got %q", expectedUser.ID, repo.findID)
	}

	if user == nil {
		t.Fatal("expected user, got nil")
	}

	if *user != expectedUser {
		t.Fatalf("expected user %+v, got %+v", expectedUser, *user)
	}
}

func TestServiceGetUserByIDReturnsRepositoryError(t *testing.T) {
	t.Parallel()

	repo := &mockRepository{
		findErr: ErrNotFound,
	}
	service := NewService(repo)

	user, err := service.GetUserByID(context.Background(), "missing-user-id")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected error %v, got %v", ErrNotFound, err)
	}

	if user != nil {
		t.Fatalf("expected nil user on error, got %+v", user)
	}
}

func TestServiceCreateUser(t *testing.T) {
	t.Parallel()

	expectedUser := newTestUser("8b7d8a0a-6c7e-4bb2-bb6a-68c73a4b6b01", "ira.bondar@example.com", "ira")
	repo := &mockRepository{
		createResult: &expectedUser,
	}
	service := NewService(repo)

	user, err := service.CreateUser(context.Background(), CreateUserInput{
		Email:    " ira.bondar@example.com ",
		Username: " ira ",
		Password: "secret123",
	})
	if err != nil {
		t.Fatalf("CreateUser returned error: %v", err)
	}

	if !repo.createCalled {
		t.Fatal("expected repository Create to be called")
	}

	if repo.createUser.Email != "ira.bondar@example.com" {
		t.Fatalf("expected repository user email to be trimmed, got %q", repo.createUser.Email)
	}

	if repo.createUser.Username != "ira" {
		t.Fatalf("expected repository user username to be trimmed, got %q", repo.createUser.Username)
	}

	if repo.createUser.PasswordHash == "" {
		t.Fatal("expected repository user password hash to be set")
	}

	if repo.createUser.PasswordHash == "secret123" {
		t.Fatal("expected repository user password to be hashed")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(repo.createUser.PasswordHash), []byte("secret123")); err != nil {
		t.Fatalf("expected repository user password hash to match password: %v", err)
	}

	if user == nil {
		t.Fatal("expected created user, got nil")
	}

	if *user != expectedUser {
		t.Fatalf("expected created user %+v, got %+v", expectedUser, *user)
	}
}

func TestServiceCreateUserRequiresEmail(t *testing.T) {
	t.Parallel()

	repo := &mockRepository{}
	service := NewService(repo)

	user, err := service.CreateUser(context.Background(), CreateUserInput{
		Email:    " ",
		Username: "ira",
		Password: "secret123",
	})
	if !errors.Is(err, ErrEmailMustPresent) {
		t.Fatalf("expected error %v, got %v", ErrEmailMustPresent, err)
	}

	if user != nil {
		t.Fatalf("expected nil user, got %+v", user)
	}

	if repo.createCalled {
		t.Fatal("expected repository Create not to be called")
	}
}

func TestServiceCreateUserRequiresUsername(t *testing.T) {
	t.Parallel()

	repo := &mockRepository{}
	service := NewService(repo)

	user, err := service.CreateUser(context.Background(), CreateUserInput{
		Email:    "ira.bondar@example.com",
		Username: " ",
		Password: "secret123",
	})
	if !errors.Is(err, ErrUsernameMustPresent) {
		t.Fatalf("expected error %v, got %v", ErrUsernameMustPresent, err)
	}

	if user != nil {
		t.Fatalf("expected nil user, got %+v", user)
	}

	if repo.createCalled {
		t.Fatal("expected repository Create not to be called")
	}
}

func TestServiceCreateUserRequiresPassword(t *testing.T) {
	t.Parallel()

	repo := &mockRepository{}
	service := NewService(repo)

	user, err := service.CreateUser(context.Background(), CreateUserInput{
		Email:    "ira.bondar@example.com",
		Username: "ira",
		Password: " ",
	})
	if !errors.Is(err, ErrPasswordMustPresent) {
		t.Fatalf("expected error %v, got %v", ErrPasswordMustPresent, err)
	}

	if user != nil {
		t.Fatalf("expected nil user, got %+v", user)
	}

	if repo.createCalled {
		t.Fatal("expected repository Create not to be called")
	}
}

func TestServiceCreateUserReturnsRepositoryError(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("create user failed")
	repo := &mockRepository{
		createErr: expectedErr,
	}
	service := NewService(repo)

	user, err := service.CreateUser(context.Background(), CreateUserInput{
		Email:    "ira.bondar@example.com",
		Username: "ira",
		Password: "secret123",
	})
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}

	if user != nil {
		t.Fatalf("expected nil user, got %+v", user)
	}
}

func newTestUser(id, email, username string) User {
	now := time.Date(2026, 8, 24, 10, 0, 0, 0, time.UTC)

	return User{
		ID:           id,
		Email:        email,
		Username:     username,
		PasswordHash: "password-hash",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
