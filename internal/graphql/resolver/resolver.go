package resolver

import "github.com/vladfc/ghira/internal/user"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	UserService *user.Service
}
