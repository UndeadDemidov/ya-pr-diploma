DROP TABLE users;

CREATE TABLE users
(
    id         UUID                      NOT NULL
        CONSTRAINT users_pk
            PRIMARY KEY,
    created_at timestamptz DEFAULT NOW() NOT NULL
);