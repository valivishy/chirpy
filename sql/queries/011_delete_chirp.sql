-- name: DeleteChirp :exec
DELETE FROM chirps where id = $1;
