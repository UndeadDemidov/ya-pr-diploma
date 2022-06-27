package postgre

import (
	"context"

	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/auth"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/user"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	selectUserByLogin = "SELECT user_id FROM auth WHERE login=$1"
	// selectUserByLogin = "SELECT user_id FROM auth WHERE login=$1 AND PASSWORD=$2"
	insertCredentials = "INSERT INTO auth (user_id, login, password) VALUES ($1, $2, $3)"
)

type Auth struct {
	db *pgxpool.Pool
}

var _ auth.Repository = (*Auth)(nil)

func NewAuth(db *pgxpool.Pool) *Auth {
	return &Auth{db: db}
}

func (a Auth) Create(ctx context.Context, usr user.User, login, pword string) error {
	_, err := a.db.Exec(ctx, insertCredentials, usr.ID, login, pword)
	if err != nil {
		return err
	}
	return nil
}

func (a Auth) Read(ctx context.Context, login string) (usr user.User, err error) {
	var userID string
	err = a.db.QueryRow(ctx, selectUserByLogin, login).Scan(&userID)
	if err != nil {
		return user.User{}, err
	}
	return user.User{ID: userID}, nil
}
