-- name: UpdateUser :exec
UPDATE "user"
SET full_name = $2,
    bio = $3,
    email = $4,
    updated_at = NOW()
WHERE id = $1;

-- name: UpdateUserPassword :exec
UPDATE "user"
SET password_hash = $2,
    updated_at = NOW()
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM "user"
WHERE email = $1;

-- name: DeactivateUser :exec
UPDATE "user"
SET is_active = FALSE,
    updated_at = NOW()
WHERE id = $1;

-- name: ActivateUser :exec
UPDATE "user"
SET is_active = TRUE,
    updated_at = NOW()
WHERE id = $1;

-- name: SearchUsersByName :many
SELECT * FROM "user"
WHERE full_name ILIKE '%' || $1 || '%'
ORDER BY full_name;

-- name: AddUser :one
INSERT INTO "user" (
  full_name,
  bio,
  email,
  password_hash,
  created_at,
  updated_at
) VALUES (
  $1, $2, $3, $4, NOW(), NOW()
)
RETURNING id;

-- name: GetUserByID :one
SELECT * FROM "user"
WHERE id = $1;