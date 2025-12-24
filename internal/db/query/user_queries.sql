-- name: CreateUser :one
INSERT INTO users (
    username,
    hashed_password,
    first_name,
    last_name,
    email,
    phone,
    role
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmailOrUsername :one
SELECT * FROM users
WHERE email = $1 OR username = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
WHERE 
    (sqlc.narg('role')::user_role IS NULL OR role = sqlc.narg('role'))
    AND (sqlc.narg('search')::text IS NULL OR 
         full_name ILIKE '%' || sqlc.narg('search') || '%' OR 
         email ILIKE '%' || sqlc.narg('search') || '%' OR 
         username ILIKE '%' || sqlc.narg('search') || '%')
    AND (
        sqlc.narg('cursor_created_at')::timestamptz IS NULL OR 
        (created_at, id) < (sqlc.narg('cursor_created_at')::timestamptz, sqlc.narg('cursor_id')::uuid)
    )
ORDER BY created_at DESC, id DESC
LIMIT $1;

-- name: UpdateUser :one
UPDATE users
SET 
    first_name = COALESCE(sqlc.narg('first_name'), first_name),
    last_name = COALESCE(sqlc.narg('last_name'), last_name),
    email = COALESCE(sqlc.narg('email'), email),
    phone = COALESCE(sqlc.narg('phone'), phone),
    role = COALESCE(sqlc.narg('role'), role),
    is_email_verified = COALESCE(sqlc.narg('is_email_verified'), is_email_verified),
    is_active = COALESCE(sqlc.narg('is_active'), is_active)
WHERE id = $1
RETURNING *;

-- name: UpdatePassword :exec
UPDATE users
SET hashed_password = $2
WHERE id = $1;

-- name: DeleteUser :exec
UPDATE users
SET deleted_at = now()
WHERE id = $1;

-- name: HardDeleteUser :exec
DELETE FROM users
WHERE id = $1;
