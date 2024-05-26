-- +goose Up

-- create books table
CREATE TABLE IF NOT EXISTS catalog.books
(
    book_id    BIGINT PRIMARY KEY                                               NOT NULL,
    book_title TEXT                                                             NOT NULL,
    book_isbn  TEXT                                                             NOT NULL,
    book_desc  TEXT                                                             NOT NULL,
    book_price DECIMAL(10, 2)                                                   NOT NULL,
    author_id  BIGINT REFERENCES catalog.authors (author_id) ON DELETE RESTRICT NOT NULL,
    deleted    BOOLEAN     DEFAULT FALSE                                        NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()                                        NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW()                                        NOT NULL
);

-- +goose Down
