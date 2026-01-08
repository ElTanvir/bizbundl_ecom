-- name: CreateShop :one
INSERT INTO shops (
    owner_id,
    name,
    subdomain,
    tenant_id,
    is_active
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetShopBySubdomain :one
SELECT * FROM shops
WHERE subdomain = $1 LIMIT 1;

-- name: ListShopsByOwner :many
SELECT * FROM shops
WHERE owner_id = $1;
