-- +goose Up
CREATE TABLE comics (
    month TEXT,
    num INT,
    link TEXT,
    year TEXT,
    news TEXT,
    safe_title TEXT,
    transcript TEXT,
    alt TEXT,
    img TEXT,
    title TEXT,
    day TEXT
);

-- +goose Down
DROP TABLE IF EXISTS comics;