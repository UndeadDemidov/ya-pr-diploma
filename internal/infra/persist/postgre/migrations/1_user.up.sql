CREATE TABLE users
(
    id         UUID                      NOT NULL
        CONSTRAINT users_pk
            PRIMARY KEY,
    balance    INTEGER     DEFAULT 0     NOT NULL,
    accrual    INTEGER     DEFAULT 0     NOT NULL,
    withdrawn  INTEGER     DEFAULT 0     NOT NULL,
    created_at timestamptz DEFAULT NOW() NOT NULL,
    updated_at timestamptz DEFAULT NOW() NOT NULL
);