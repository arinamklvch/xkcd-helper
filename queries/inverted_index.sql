-- name: InsertIntoInvertedIndex :copyfrom
INSERT INTO inverted_index(
    word,
    comics_nums
) VALUES (
    $1, $2
);

-- name: GetFromInvertedIndex :many
SELECT comics_nums
FROM inverted_index
WHERE word = ANY(@words::TEXT[]);   