// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: 005_delete_chirps.sql

package database

import (
	"context"
)

const deleteChirps = `-- name: DeleteChirps :exec
DELETE FROM chirps
`

func (q *Queries) DeleteChirps(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteChirps)
	return err
}
