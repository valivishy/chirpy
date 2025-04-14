// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: 001_users.sql

package database

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (email, hashed_password)
VALUES (
   $1,
        $2
)
RETURNING id, created_at, updated_at, email, hashed_password
`

type CreateUserParams struct {
	Email          string
	HashedPassword string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Email, arg.HashedPassword)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
	)
	return i, err
}
