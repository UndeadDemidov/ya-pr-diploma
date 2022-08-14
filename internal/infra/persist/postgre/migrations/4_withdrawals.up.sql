CREATE TABLE withdrawals
(
    id           UUID        DEFAULT gen_random_uuid() NOT NULL
        CONSTRAINT withdrawals_pk
            PRIMARY KEY,
    user_id      uuid                                  NOT NULL
        CONSTRAINT withdrawals_user_id_fk
            REFERENCES users,
    order_number BIGINT                                NOT NULL,
    sum          INTEGER     DEFAULT 0                 NOT NULL,
    processed_at timestamptz DEFAULT NOW()             NOT NULL
);

CREATE OR REPLACE FUNCTION update_withdrawn()
    RETURNS TRIGGER AS
$$
BEGIN
    IF NEW.sum > (SELECT u.balance FROM users u WHERE u.id = NEW.user_id) THEN
        RAISE EXCEPTION 'cannot write off sum more than balance';
    END IF;

    UPDATE users AS u
    SET withdrawn  = u.withdrawn + NEW.sum,
        balance    = u.balance - NEW.sum,
        updated_at = NOW()
    WHERE u.id = NEW.user_id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_before_insert_withdrawals ON withdrawals;

CREATE TRIGGER trigger_before_insert_withdrawals
    BEFORE INSERT
    ON withdrawals
    FOR EACH ROW
EXECUTE PROCEDURE update_withdrawn();