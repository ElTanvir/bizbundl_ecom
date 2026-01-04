-- name: CreateStoreConfig :one
INSERT INTO store_configs (
    key,
    value,
    is_encrypted,
    group_name
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetStoreConfig :one
SELECT * FROM store_configs
WHERE key = $1 LIMIT 1;

-- name: ListStoreConfigs :many
SELECT * FROM store_configs;

-- name: UpdateStoreConfig :one
UPDATE store_configs
SET 
    value = $2,
    is_encrypted = $3,
    group_name = $4,
    updated_at = now()
WHERE key = $1
RETURNING *;

-- name: DeleteStoreConfig :exec
DELETE FROM store_configs
WHERE key = $1;
