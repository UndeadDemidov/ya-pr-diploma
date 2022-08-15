package postgre

import (
	"context"
	"embed"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
)

//go:embed migrations/*.sql
var fs embed.FS

type Persist struct {
	*User
	*Auth
	*Order
	*Balance
}

func NewPersist(ctx context.Context, db *pgxpool.Pool) (*Persist, error) {
	if db == nil {
		panic("missing *pgxpool.Pool, parameter must not be nil")
	}
	err := db.Ping(ctx)
	if err != nil {
		return nil, err
	}
	// проверить что бд еще не создано и создать
	log.Info().
		Str("database", db.Config().ConnConfig.Database).
		Str("host", db.Config().ConnConfig.Host).
		Msg("successfully connected to PG")

	err = migrateDB(db)
	if err != nil {
		return nil, err
	}

	return &Persist{
		User:    NewUser(db),
		Auth:    NewAuth(db),
		Order:   NewOrder(db),
		Balance: NewBalance(db),
	}, nil
}

// migrateDB хотел сделать через golang-migrate/migrate - но только потерял время.
// несовместимые connection string и нельзя конвертировать нативный постгресовый формат в uri
func migrateDB(db *pgxpool.Pool) error {
	d, err := iofs.New(fs, "migrations")
	if err != nil {
		return err
	}

	log.Debug().Msg("starting migrations")
	m, err := migrate.NewWithSourceInstance("iofs", d, db.Config().ConnConfig.ConnString())
	if err != nil {
		return err
	}

	log.Debug().Msg("starting upgrades")
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	log.Info().Msg("DB is initialized successfully")
	return nil
}
