package postgre

import (
	"context"
	"errors"

	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/order"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/primit"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/user"
	errors2 "github.com/UndeadDemidov/ya-pr-diploma/internal/errors"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Order struct {
	db *pgxpool.Pool
}

var _ order.Repository = (*Order)(nil)

func NewOrder(db *pgxpool.Pool) *Order {
	if db == nil {
		panic("missing *pgxpool.Pool, parameter must not be nil")
	}
	return &Order{db: db}
}

func (o Order) Create(ctx context.Context, ord order.Order) error {
	const insertQuery = `
WITH inserted_rows AS (
    INSERT INTO orders (id, user_id, number)
        VALUES ($1, $2, $3)
        ON CONFLICT (number) DO NOTHING
        RETURNING id)
SELECT user_id
FROM orders
WHERE NOT EXISTS(SELECT 1 FROM inserted_rows)
  AND number = $3;
`
	var usrID string
	err := o.db.QueryRow(ctx, insertQuery, ord.ID, ord.User.ID, uint64(ord.Number)).Scan(&usrID)
	if errors.Is(err, pgx.ErrNoRows) {
		// Если пустой сет записей, то успешно вставили запись
		return nil
	}
	if err != nil {
		// Если другая ошибка - возвращаем ее
		return err
	}
	// Если ошибок нет, значит нашли в БД такой заказ
	if usrID == ord.User.ID {
		// Если вернулся ID пользователя, значит заказ уже загружен
		return errors2.ErrOrderAlreadyUploaded
	}
	// Если вернулся ID другого пользователя, значит заказ уже загружен другим пользователем
	return errors2.ErrOrderAlreadyUploadedByAnotherUser
}

func (o Order) ListByUser(ctx context.Context, usr user.User) ([]order.Order, error) {
	const selectOrdersByUser = `
SELECT id, number, status, accrual, uploaded_at, processed_at
FROM orders
WHERE user_id=$1
`
	rows, err := o.db.Query(ctx, selectOrdersByUser, usr.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ords := make([]order.Order, 0, 4)
	for rows.Next() {
		var (
			status string
			money  int64
		)
		ord := order.Order{}
		err = rows.Scan(&ord.ID, &ord.Number, &status, &money, &ord.Unloaded, &ord.Processed)
		if err != nil {
			return nil, err
		}
		procStatus, err := order.ParseProcessingStatus(status)
		if err != nil {
			return nil, err
		}
		ord.Status = procStatus
		ord.Accrual = primit.Currency(money)
		ords = append(ords, ord)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return ords, nil
}
