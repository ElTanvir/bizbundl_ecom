-- name: CreateCart :one
INSERT INTO carts (
    session_id,
    user_id,
    status
) VALUES (
    $1, $2, 'active'
) RETURNING *;

-- name: GetCartBySession :one
SELECT * FROM carts
WHERE session_id = $1 AND status = 'active'
LIMIT 1;

-- name: GetCartByUser :one
SELECT * FROM carts
WHERE user_id = $1 AND status = 'active'
LIMIT 1;

-- name: UpdateCartUser :exec
UPDATE carts
SET user_id = $2
WHERE id = $1;

-- name: AddCartItem :one
INSERT INTO cart_items (
    cart_id,
    product_id,
    variant_id,
    quantity
) VALUES (
    $1, $2, $3, $4
) 
ON CONFLICT (cart_id, product_id, COALESCE(variant_id, '00000000-0000-0000-0000-000000000000'::uuid)) 
DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity
RETURNING *;

-- name: UpdateCartItemQuantity :one
UPDATE cart_items
SET quantity = $3
WHERE id = $1 AND cart_id = $2
RETURNING *;

-- name: RemoveCartItem :exec
DELETE FROM cart_items
WHERE id = $1 AND cart_id = $2;

-- name: GetCartItems :many
SELECT ci.*, p.title as product_title, p.base_price, p.is_digital, p.file_path,
       pv.title as variant_title, pv.price as variant_price
FROM cart_items ci
JOIN products p ON ci.product_id = p.id
LEFT JOIN product_variants pv ON ci.variant_id = pv.id
WHERE ci.cart_id = $1
ORDER BY ci.created_at ASC;

-- name: ClearCart :exec
DELETE FROM cart_items
WHERE cart_id = $1;

-- name: DeleteCart :exec
DELETE FROM carts
WHERE id = $1;
