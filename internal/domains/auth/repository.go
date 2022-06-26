package auth

import (
	"context"

	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/user"
)

// Repository представляет собой классический CRUD
type Repository interface {
	// Fetch(ctx context.Context) ([]*User, error)
	Validate(ctx context.Context, user *user.User, login, pword string) bool
	Create(ctx context.Context, user *user.User) error
	// Update(ctx context.Context, user *User) (*User, error)
	// Delete(ctx context.Context, id int) error
}
