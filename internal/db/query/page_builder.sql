-- name: CreatePage :one
INSERT INTO pages (
    route,
    name,
    sections,
    is_published
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetPageByRoute :one
SELECT * FROM pages
WHERE route = $1 LIMIT 1;

-- name: ListPages :many
SELECT * FROM pages
ORDER BY updated_at DESC;

-- name: UpdatePage :one
UPDATE pages
SET 
    sections = $2,
    is_published = $3,
    updated_at = NOW()
WHERE route = $1
RETURNING *;
