-- +goose Up
CREATE TABLE inverted_index (
    word TEXT,
    comics_nums INT[]
);

-- +goose Down
DROP TABLE IF EXISTS inverted_index;