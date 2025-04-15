-- name: GetUserFromRefreshToken :one
SELECT u.*
FROM users u
INNER JOIN refresh_tokens rf on rf.user_id = u.id
WHERE rf.token = $1
  AND rf.revoked_at IS NULL
  AND rf.expires_at >= NOW()
LIMIT 1;