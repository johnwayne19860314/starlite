-- name: CreateUser :one
INSERT INTO first.user (
  user_name,
  user_email,
  user_password,
  user_role,
  is_active
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetUser :one
SELECT * FROM first.user
WHERE id = $1 LIMIT 1;

-- name: GetUserByName :one
SELECT * FROM first.user
WHERE user_name = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: ListUsers :many
SELECT * FROM first.user
WHERE is_active = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: UpdateUserRole :one
UPDATE first.user
SET user_role = $2
WHERE id = $1
RETURNING *;


-- name: UpdateUserPassword :one
UPDATE first.user
SET user_password = $2
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
UPDATE first.user
SET is_active = false
WHERE id = $1
RETURNING *;
