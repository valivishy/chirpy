// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: 012_update_user_chirpy_red.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const updateUserChirpyRed = `-- name: UpdateUserChirpyRed :exec
UPDATE users
SET is_chirpy_red = true WHERE id = $1
`

func (q *Queries) UpdateUserChirpyRed(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, updateUserChirpyRed, id)
	return err
}
