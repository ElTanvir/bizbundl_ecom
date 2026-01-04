-- name: CreatePaymentGateway :one
INSERT INTO payment_gateways (
    id,
    name,
    config,
    is_test_mode,
    is_active,
    position
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetPaymentGateway :one
SELECT * FROM payment_gateways
WHERE id = $1 LIMIT 1;

-- name: ListPaymentGateways :many
SELECT * FROM payment_gateways
ORDER BY position ASC;

-- name: UpdatePaymentGateway :one
UPDATE payment_gateways
SET 
    name = COALESCE(sqlc.narg('name'), name),
    config = COALESCE(sqlc.narg('config'), config),
    is_test_mode = COALESCE(sqlc.narg('is_test_mode'), is_test_mode),
    is_active = COALESCE(sqlc.narg('is_active'), is_active),
    position = COALESCE(sqlc.narg('position'), position)
WHERE id = $1
RETURNING *;

-- name: DeletePaymentGateway :exec
DELETE FROM payment_gateways
WHERE id = $1;
