-- name: GetUser :one
SELECT *
FROM users
WHERE login = sqlc.arg(login)
  AND password = sqlc.arg(password);