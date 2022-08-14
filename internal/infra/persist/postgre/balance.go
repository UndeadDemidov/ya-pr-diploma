package postgre

import (
	"context"

	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/balance"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/primit"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/user"
	errors2 "github.com/UndeadDemidov/ya-pr-diploma/internal/errors"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
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

func (b Balance) CreateWithdrawal(ctx context.Context, wtdrwl balance.Withdrawal) error {
	const insertWtdrwl = `
INSERT INTO withdrawals (id, user_id, order_number, sum)
VALUES ($1, $2, $3, $4);
`
	tx, err := b.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			err = tx.Rollback(ctx)
			if err != nil {
				log.Err(err).Msg("got error on rollback transaction")
				return
			}
		} else {
			err = tx.Commit(ctx)
			if err != nil {
				log.Err(err).Msg("got error on commit transaction")
				return
			}
		}
	}()

	bal, err := b.Read(ctx, wtdrwl.User)
	if err != nil {
		return err
	}
	if wtdrwl.Sum > bal.Current {
		return errors2.ErrWithdrawalNotEnoughFund
	}

	_, err = b.db.Exec(ctx, insertWtdrwl, wtdrwl.ID, wtdrwl.User.ID, wtdrwl.Order.String(), wtdrwl.Sum)
	if err != nil {
		return err
	}
	return nil
}

func (b Balance) ListWithdrawals(ctx context.Context, usr user.User) ([]balance.Withdrawal, error) {
	const selectWithdrawalsByUser = `
SELECT id, order_number, sum, processed_at
FROM withdrawals
WHERE user_id=$1
`
	rows, err := b.db.Query(ctx, selectWithdrawalsByUser, usr.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return b.list(usr, rows)
}

func (b Balance) list(usr user.User, rows pgx.Rows) (wtdrwls []balance.Withdrawal, err error) {
	wtdrwls = make([]balance.Withdrawal, 0, 4)
	for rows.Next() {
		var (
			money int64
		)
		wtdrwl := balance.Withdrawal{User: usr}
		err = rows.Scan(&wtdrwl.ID, &wtdrwl.Order, &money, &wtdrwl.Processed)
		if err != nil {
			return nil, err
		}
		wtdrwl.Sum = primit.Currency(money)
		wtdrwls = append(wtdrwls, wtdrwl)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return wtdrwls, nil
}
