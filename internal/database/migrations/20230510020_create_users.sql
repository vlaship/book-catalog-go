-- +goose Up

-- create users table
CREATE TABLE IF NOT EXISTS catalog.users
(
    user_id    BIGINT PRIMARY KEY NOT NULL,
    username   TEXT UNIQUE        NOT NULL,
    password   TEXT               NOT NULL,
    user_data  JSONB              NOT NULL,
    deleted    BOOLEAN            NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ                 DEFAULT NOW()
);

-- +goose Down
