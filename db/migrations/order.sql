CREATE TYPE order_status AS ENUM ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');

CREATE TABLE orders
(
    id           UUID         DEFAULT gen_random_uuid() NOT NULL
        CONSTRAINT orders_pk
            PRIMARY KEY,
    user_id      uuid                                   NOT NULL
        CONSTRAINT auth_users_id_fk
            REFERENCES users,
    number       VARCHAR                                NOT NULL,
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