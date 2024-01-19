-- name: CreateEntry :one
INSERT INTO first.entry (
  entry_code,
  entry_category,
  entry_name,
  entry_amount,
  entry_weight,
  entry_note,
  supplier_name,
  supplier_contact_info,
  fix,
  chemical_name,
  is_active
) VALUES (
  $1, $2, $3, $4, $5, $6,$7,$8, $9,$10,$11
) RETURNING *;

-- name: GetEntry :one
SELECT * FROM first.entry
WHERE id = $1 LIMIT 1;

-- name: GetEntryByName :one
SELECT * FROM first.entry
WHERE entry_name = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: GetEntryByCode :one
SELECT * FROM first.entry
WHERE entry_code = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: ListEntrys :many
SELECT * FROM first.entry
WHERE is_active = $1
and entry_category = $4
ORDER BY updated_at desc
LIMIT $2
OFFSET $3;

-- name: UpdateEntryAmount :one
UPDATE first.entry
SET entry_amount = $2
WHERE id = $1
RETURNING *;

-- name: UpdateEntryWeight :one
UPDATE first.entry
SET entry_weight = $2
WHERE id = $1
RETURNING *;

-- name: UpdateEntryNote :one
UPDATE first.entry
SET entry_note = $2
WHERE id = $1
RETURNING *;

-- name: DeleteEntry :exec
UPDATE first.entry
SET is_active = false
WHERE entry_code = $1
RETURNING *;

-- name: CreateMultipleEntries :copyfrom
INSERT INTO first.entry (
  entry_code,
  entry_category,
  entry_name,
  entry_amount,
  entry_weight,
  entry_note,
  supplier_name,
  supplier_contact_info,
  fix,
  chemical_name,
  is_active
) VALUES (
  $1, $2, $3, $4, $5, $6,$7,$8, $9,$10,$11
);