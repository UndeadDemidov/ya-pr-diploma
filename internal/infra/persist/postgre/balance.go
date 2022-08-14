package postgre

import (
	"context"

	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/balance"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/user"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Balance struct {
	db *pgxpool.Pool
}

var _ balance.Repository = (*Balance)(nil)

func NewBalance(db *pgxpool.Pool) *Balance {
	if db == nil {
		panic("missing *pgxpool.Pool, parameter must not be nil")
	}
	return &Balance{db: db}
}

func (b Balance) Read(ctx context.Context, usr user.User) (balance.Balance, error) {
	query := `
select balance, accrual, withdrawn
from users
where id=$1;
`
	var bal, acc, wth int64
	err := b.db.QueryRow(ctx, query, usr.ID).Scan(&bal, &acc, &wth)
	if err != nil {
		return balance.Balance{}, err
	}
	return balance.NewBalance(usr, bal, acc, wth), nil
}
