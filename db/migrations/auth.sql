DROP TABLE auth;

CREATE TABLE auth
(
    id       UUID DEFAULT gen_random_uuid() NOT NULL
        CONSTRAINT auth_pk
            PRIMARY KEY,
    user_id  uuid                           NOT NULL
        CONSTRAINT auth_users_id_fk
            REFERENCES users,
    login    VARCHAR                        NOT NULL,
    password VARCHAR                        NOT NULL
);

CREATE UNIQUE INDEX auth_login_uindex
    ON auth (login);

CREATE UNIQUE INDEX auth_user_id_uindex
    ON auth (user_id);
