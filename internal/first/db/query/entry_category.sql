-- name: CreateEntryCategory :one
INSERT INTO first.entry_category (
  category,
  note,
  is_active
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetEntryCategory :one
SELECT * FROM first.entry_category
WHERE category = $1 LIMIT 1;

-- name: ListEntryCategories :many
SELECT * FROM first.entry_category
WHERE is_active = $1
ORDER BY id;


-- name: DeleteEntryCategory :exec
UPDATE first.entry_category
SET is_active = false
WHERE category = $1
RETURNING *;


-- name: UpdateEntryCategory :exec
UPDATE first.entry_category
SET note = $1
WHERE category = $2
RETURNING *;
