package user

import (
	"context"
	"database/sql"
	"errors"

	userdb "github.com/vladfc/ghira/internal/database/db"
)

type PostgresRepository struct {
	queries *userdb.Queries
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		queries: userdb.New(db),
	}
}

func (r *PostgresRepository) Create(ctx context.Context, user User) (*User, error) {
	row, err := r.queries.CreateUser(ctx, toCreateUserParams(user))
	if err != nil {
		return nil, err
	}

	createdUser := toUser(row)
	return &createdUser, nil
}

func (r *PostgresRepository) List(ctx context.Context) ([]User, error) {
	rows, err := r.queries.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	return toUsers(rows), nil
}

func (r *PostgresRepository) FindByID(ctx context.Context, id string) (*User, error) {
	row, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	foundUser := toUser(row)
	return &foundUser, nil
}

func toCreateUserParams(user User) userdb.CreateUserParams {
	return userdb.CreateUserParams{
		Email:        user.Email,
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
	}
}

func toUser(row userdb.User) User {
	return User{
		ID:           row.ID,
		Email:        row.Email,
		Username:     row.Username,
		PasswordHash: row.PasswordHash,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
}

func toUsers(rows []userdb.User) []User {
	users := make([]User, 0, len(rows))

	for _, row := range rows {
		users = append(users, toUser(row))
	}

	return users
}

var _ Repository = (*PostgresRepository)(nil)
