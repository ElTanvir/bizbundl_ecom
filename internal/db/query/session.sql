-- name: CreateSession :one
INSERT INTO sessions (
    token,
    user_id,
    expires_at
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetSession :one
SELECT * FROM sessions
WHERE token = $1 LIMIT 1;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE token = $1;

-- name: DeleteExpiredSessions :exec
DELETE FROM sessions
WHERE expires_at < now();
