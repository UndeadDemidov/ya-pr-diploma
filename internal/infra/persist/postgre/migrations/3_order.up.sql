CREATE TYPE order_status AS ENUM ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');

CREATE TABLE orders
(
    id           UUID         DEFAULT gen_random_uuid() NOT NULL
        CONSTRAINT orders_pk
            PRIMARY KEY,
    user_id      uuid                                   NOT NULL
        CONSTRAINT orders_users_id_fk
            REFERENCES users,
    number       BIGINT                                 NOT NULL,
    status       order_status DEFAULT 'NEW'             NOT NULL,
    accrual      INTEGER      DEFAULT 0 				NOT NULL,
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

CREATE OR REPLACE FUNCTION update_accrual()
    RETURNS TRIGGER AS
$$
BEGIN
    UPDATE users as u
    SET accrual = u.accrual + (NEW.accrual - OLD.accrual),
        balance = u.balance + (NEW.accrual - OLD.accrual),
        updated_at = now()
    WHERE u.id = NEW.user_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_after_update_orders
    AFTER UPDATE
    ON orders
    FOR EACH ROW
    WHEN (OLD.status IS DISTINCT FROM NEW.status
        AND NEW.status = 'PROCESSED'
        AND OLD.accrual IS DISTINCT FROM NEW.accrual)
EXECUTE PROCEDURE update_accrual();