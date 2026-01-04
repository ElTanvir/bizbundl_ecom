-- Categories

-- name: CreateCategory :one
INSERT INTO categories (
    name,
    slug,
    parent_id,
    is_active
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetCategory :one
SELECT * FROM categories
WHERE id = $1 LIMIT 1;

-- name: ListCategories :many
SELECT * FROM categories
ORDER BY name ASC;

-- name: UpdateCategory :one
UPDATE categories
SET 
    name = COALESCE(sqlc.narg('name'), name),
    slug = COALESCE(sqlc.narg('slug'), slug),
    parent_id = COALESCE(sqlc.narg('parent_id'), parent_id),
    is_active = COALESCE(sqlc.narg('is_active'), is_active)
WHERE id = $1
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories
WHERE id = $1;


-- Products

-- name: CreateProduct :one
INSERT INTO products (
    title,
    slug,
    description,
    base_price,
    is_digital,
    file_path,
    category_id,
    is_active
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetProduct :one
SELECT * FROM products
WHERE id = $1 LIMIT 1;

-- name: GetProductBySlug :one
SELECT * FROM products
WHERE slug = $1 LIMIT 1;

-- name: ListProducts :many
SELECT * FROM products
ORDER BY created_at DESC;

-- name: UpdateProduct :one
UPDATE products
SET 
    title = COALESCE(sqlc.narg('title'), title),
    slug = COALESCE(sqlc.narg('slug'), slug),
    description = COALESCE(sqlc.narg('description'), description),
    base_price = COALESCE(sqlc.narg('base_price'), base_price),
    is_digital = COALESCE(sqlc.narg('is_digital'), is_digital),
    file_path = COALESCE(sqlc.narg('file_path'), file_path),
    category_id = COALESCE(sqlc.narg('category_id'), category_id),
    is_active = COALESCE(sqlc.narg('is_active'), is_active)
WHERE id = $1
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1;


-- Variants

-- name: CreateProductVariant :one
INSERT INTO product_variants (
    product_id,
    title,
    options,
    price,
    compare_at_price,
    sku,
    stock_quantity,
    is_active
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: ListVariantsByProduct :many
SELECT * FROM product_variants
WHERE product_id = $1;

-- name: DeleteVariant :exec
DELETE FROM product_variants
WHERE id = $1;
