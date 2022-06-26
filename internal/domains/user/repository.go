package user

import "context"

type Repository interface {
	Fetch(ctx context.Context) ([]*User, error)
	FetchByID(ctx context.Context, id int) (*User, error)
	Create(ctx context.Context, user *User) (*User, error)
	// Update(ctx context.Context, user *User) (*User, error)
	// Delete(ctx context.Context, id int) error
}
