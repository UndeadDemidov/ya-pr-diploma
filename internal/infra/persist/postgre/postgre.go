package postgre

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Persist struct {
	*User
	*Auth
}

func NewPersist(ctx context.Context, db *pgxpool.Pool) (*Persist, error) {
	err := db.Ping(ctx)
	if err != nil {
		return nil, err
	}
	// проверить что бд еще не создано и создать
	return &Persist{
		User: NewUser(db),
		Auth: NewAuth(db),
	}, nil
}
