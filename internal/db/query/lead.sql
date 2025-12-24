-- name: CreateLead :exec
INSERT INTO leads (
  id,
  first_name,
  last_name,
  email,
  phone,
  company,
  message,
  status,
  source
) VALUES (
  ?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8, ?9
);


-- name: ListLeadsFiltered :many
SELECT
  id,
  first_name,
  last_name,
  email,
  phone,
  company,
  message,
  status,
  source,
  created_at
FROM leads
WHERE
  (?1 IS NULL OR status = ?2)
  AND (?3 IS NULL OR created_at >= ?4)
  AND (?5 IS NULL OR created_at <= ?6)
  AND (
    ?7 IS NULL
    OR lower(first_name) LIKE lower('%' || ?8 || '%')
    OR lower(last_name) LIKE lower('%' || ?9 || '%')
    OR lower(email) LIKE lower('%' || ?10 || '%')
  )
ORDER BY created_at DESC
LIMIT ?11 OFFSET ?12;
