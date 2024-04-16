// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: query.sql

package db

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
  name, email, password
) VALUES (
  $1, $2, $3
)
RETURNING id, name, password, email
`

type CreateUserParams struct {
	Name     string
	Email    string
	Password string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser, arg.Name, arg.Email, arg.Password)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Password,
		&i.Email,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT id, name, password, email FROM users
WHERE name = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, name string) (User, error) {
	row := q.db.QueryRow(ctx, getUser, name)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Password,
		&i.Email,
	)
	return i, err
}

const updatePassword = `-- name: UpdatePassword :exec
UPDATE users
  set password = $2
WHERE id = $1
`

type UpdatePasswordParams struct {
	ID       int64
	Password string
}

func (q *Queries) UpdatePassword(ctx context.Context, arg UpdatePasswordParams) error {
	_, err := q.db.Exec(ctx, updatePassword, arg.ID, arg.Password)
	return err
}

const updateUser = `-- name: UpdateUser :exec
UPDATE users
  set name = $2,
  email = $3
WHERE id = $1
`

type UpdateUserParams struct {
	ID    int64
	Name  string
	Email string
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.Exec(ctx, updateUser, arg.ID, arg.Name, arg.Email)
	return err
}
