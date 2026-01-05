-- name: CreateOrder :one
INSERT INTO orders (
    user_id,
    total_amount,
    status,
    payment_status,
    payment_method
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: CreateOrderItem :one
INSERT INTO order_items (
    order_id,
    product_id,
    variation_id,
    quantity,
    price_at_booking,
    title
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetOrder :one
SELECT * FROM orders
WHERE id = $1 LIMIT 1;

-- name: GetOrderItems :many
SELECT * FROM order_items
WHERE order_id = $1;

-- name: ListOrdersByUser :many
SELECT * FROM orders
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateOrderStatus :one
UPDATE orders
SET status = $2, payment_status = $3, updated_at = NOW()
WHERE id = $1
RETURNING *;
