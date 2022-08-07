package postgre

import (
	"context"

	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/auth"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/user"
	errors2 "github.com/UndeadDemidov/ya-pr-diploma/internal/errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

type Auth struct {
	db *pgxpool.Pool
}

var _ auth.Repository = (*Auth)(nil)

func NewAuth(db *pgxpool.Pool) *Auth {
	if db == nil {
		panic("missing *pgxpool.Pool, parameter must not be nil")
	}
	return &Auth{db: db}
}

func (a Auth) Create(ctx context.Context, usr user.User, login, pword string) error {
	const insertCredentials = "INSERT INTO auth (user_id, login, password) VALUES ($1, $2, $3)"
	_, err := a.db.Exec(ctx, insertCredentials, usr.ID, login, pword)
	if err != nil {
		var pgErr pq.Error
		if errors.As(err, &pgErr); pgErr.Code == pgerrcode.UniqueViolation {
			return errors2.ErrLoginIsInUseAlready
		}
		return err
	}
	return nil
}

func (a Auth) Read(ctx context.Context, login string) (usr user.User, err error) {
	var userID string
	const selectUserByLogin = "SELECT user_id FROM auth WHERE login=$1"
	err = a.db.QueryRow(ctx, selectUserByLogin, login).Scan(&userID)
	if err != nil {
		return user.User{}, err
	}
	return user.User{ID: userID}, nil
}

func (a Auth) ReadWithPassword(ctx context.Context, login, pword string) (usr user.User, err error) {
	var userID string
	const selectAuthentication = "SELECT user_id FROM auth WHERE login=$1 AND PASSWORD=$2"
	err = a.db.QueryRow(ctx, selectAuthentication, login, pword).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user.User{}, errors2.ErrPairLoginPwordIsNotExist
		}
		return user.User{}, err
	}
	return user.User{ID: userID}, nil
}
