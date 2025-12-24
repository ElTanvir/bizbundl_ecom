-- name: GetCSSVariableByName :one
SELECT * FROM css_variables WHERE name = ?;

-- name: GetAllCSSVariables :many
SELECT * FROM css_variables ORDER BY variable_type, name;

-- name: GetCSSVariablesByType :many
SELECT * FROM css_variables WHERE variable_type = ? ORDER BY name;

-- name: CreateCSSVariable :exec
INSERT INTO css_variables (name, value, variable_type)
VALUES (?, ?, ?);

-- name: UpdateCSSVariable :exec
UPDATE css_variables
SET value = ?, variable_type = ?, updated_at = CURRENT_TIMESTAMP
WHERE name = ?;

-- name: DeleteCSSVariable :exec
DELETE FROM css_variables WHERE name = ? AND is_system = FALSE;
