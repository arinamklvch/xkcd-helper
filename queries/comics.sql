-- name: InsertComics :copyfrom
INSERT INTO comics(
    month,
    num,
    link,
    year,
    news,
    safe_title,
    transcript,
    alt,
    img,
    title,
    day
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
);

-- name: GetLatestComicNum :one
SELECT COALESCE(MAX(num), 0)::INT
FROM comics;

-- name: GetComicsRange :many
SELECT *
FROM comics
WHERE sqlc.arg('from') <= num AND num <= sqlc.arg('to');

-- name: GetComicsByNums :many
SELECT *
FROM comics
WHERE num = ANY(@nums::INT[]);