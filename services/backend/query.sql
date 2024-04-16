-- name: GetUser :one
SELECT * FROM users
WHERE name = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
  name, email, password
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: UpdateUser :exec
UPDATE users
  set name = $2,
  email = $3
WHERE id = $1;

-- name: UpdatePassword :exec
UPDATE users
  set password = $2
WHERE id = $1;
