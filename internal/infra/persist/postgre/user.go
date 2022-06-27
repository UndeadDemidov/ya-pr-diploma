package postgre

import (
	"context"

	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/user"
	"github.com/jackc/pgx/v4/pgxpool"
)

const insertUser = "INSERT INTO users (id) VALUES ($1)"

type User struct {
	db *pgxpool.Pool
}

var _ user.Repository = (*User)(nil)

func NewUser(db *pgxpool.Pool) *User {
	return &User{db: db}
}

func (u User) Create(ctx context.Context, usr user.User) error {
	_, err := u.db.Exec(ctx, insertUser, usr.ID)
	if err != nil {
		return err
	}
	return nil
}
