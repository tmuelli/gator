-- name: GetUserByName :one
SELECT * FROM users WHERE name LIKE $1;