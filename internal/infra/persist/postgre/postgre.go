package postgre

import (
	"context"

	_ "github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
)

type Persist struct {
	*User
	*Auth
	*Order
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
	log.Info().Msgf("successfully connected to PG server %s", db.Config().ConnConfig.Host)

	err = migrateDB(db)
	if err != nil {
		return nil, err
	}

	return &Persist{
		User:  NewUser(db),
		Auth:  NewAuth(db),
		Order: NewOrder(db),
	}, nil
}

// migrateDB хотел сделать через golang-migrate/migrate - но только потерял время.
// несовместимые connection string и нельзя конвертировать нативный постгресовый формат в uri
func migrateDB(db *pgxpool.Pool) error {
	script := `
DROP TRIGGER IF EXISTS set_timestamp ON orders;
DROP FUNCTION IF EXISTS trigger_set_timestamp;

DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS auth;
DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS order_status;

CREATE TABLE users
(
    id         UUID                      NOT NULL
        CONSTRAINT users_pk
            PRIMARY KEY,
    created_at timestamptz DEFAULT NOW() NOT NULL
);

CREATE TABLE auth
(
    id         UUID        DEFAULT gen_random_uuid() NOT NULL
        CONSTRAINT auth_pk
            PRIMARY KEY,
    user_id    uuid                                  NOT NULL
        CONSTRAINT auth_users_id_fk
            REFERENCES users,
    login      VARCHAR                               NOT NULL,
    password   VARCHAR                               NOT NULL,
    created_at timestamptz DEFAULT NOW()             NOT NULL
);

CREATE UNIQUE INDEX auth_login_uindex
    ON auth (login);

CREATE UNIQUE INDEX auth_user_id_uindex
    ON auth (user_id);

CREATE TYPE order_status AS ENUM ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');

CREATE TABLE orders
(
    id           UUID         DEFAULT gen_random_uuid() NOT NULL
        CONSTRAINT orders_pk
            PRIMARY KEY,
    user_id      uuid                                   NOT NULL
        CONSTRAINT auth_users_id_fk
            REFERENCES users,
    number       BIGINT                                 NOT NULL,
    status       order_status DEFAULT 'NEW'             NOT NULL,
    accrual      INTEGER      DEFAULT 0,
    uploaded_at  timestamptz  DEFAULT NOW()             NOT NULL,
    processed_at timestamptz  DEFAULT NOW()             NOT NULL
);

CREATE UNIQUE INDEX orders_number_uindex
    ON orders (number);

CREATE FUNCTION trigger_set_timestamp()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.processed_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_timestamp
    BEFORE
        UPDATE
    ON orders
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

`
	_, err := db.Exec(context.Background(), script)
	if err != nil {
		return err
	}
	log.Info().Msg("DB is initialized successfully")
	return nil
}
