-- +goose Up
CREATE TABLE users (
    login TEXT,
    password TEXT,
    role TEXT
);

-- +goose Down
DROP TABLE IF EXISTS users;