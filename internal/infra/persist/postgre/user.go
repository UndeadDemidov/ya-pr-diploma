package postgre

import (
	"context"

	_ "github.com/lib/pq"

	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/user"
	"github.com/jackc/pgx/v4/pgxpool"
)

type User struct {
	db *pgxpool.Pool
}

var _ user.Repository = (*User)(nil)

func NewUser(db *pgxpool.Pool) *User {
	if db == nil {
		panic("missing *pgxpool.Pool, parameter must not be nil")
	}
	return &User{db: db}
}

func (u User) Create(ctx context.Context, usr user.User) error {
	const insertUser = "INSERT INTO users (id) VALUES ($1)"
	_, err := u.db.Exec(ctx, insertUser, usr.ID)
	if err != nil {
		return err
	}
	return nil
}
