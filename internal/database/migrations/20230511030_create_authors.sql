-- +goose Up

-- create authors table
CREATE TABLE IF NOT EXISTS catalog.authors
(
    author_id   BIGINT PRIMARY KEY        NOT NULL,
    author_name TEXT                      NOT NULL,
    author_dob  DATE                      NOT NULL,
    deleted     BOOLEAN     DEFAULT FALSE NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at  TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

-- +goose Down
