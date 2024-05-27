-- +goose Up

-- create properties table
create table catalog.properties
(
    property_id    SERIAL PRIMARY KEY,
    property_name  TEXT UNIQUE NOT NULL,
    property_value JSONB       NOT NULL
);
