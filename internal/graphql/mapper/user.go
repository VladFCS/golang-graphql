package mapper

import (
	"github.com/vladfc/ghira/internal/graphql/model"
	"github.com/vladfc/ghira/internal/user"
)

func UserToGraphQL(u user.User) *model.User {
	return &model.User{
		ID:        u.ID,
		Email:     u.Email,
		Username:  u.Username,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func UsersToGraphQL(users []user.User) []*model.User {
	result := make([]*model.User, 0, len(users))

	for _, u := range users {
		result = append(result, UserToGraphQL(u))
	}

	return result
}
