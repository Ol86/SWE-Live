-- name: GetMemberByID :one
SELECT
    id,
    version,
    username,
    first_name,
    last_name,
    gender,
    date_of_birth,
    member_since,
    is_student,
    email_address,
    interests,
    generated,
    updated
FROM library.member
WHERE id = $1;

-- name: SearchMembers :many
SELECT
    id,
    version,
    username,
    first_name,
    last_name,
    gender,
    date_of_birth,
    member_since,
    is_student,
    email_address,
    interests,
    generated,
    updated
FROM library.member
WHERE
    (sqlc.narg('username')::text IS NULL OR username ILIKE '%' || sqlc.narg('username')::text || '%')
    AND (sqlc.narg('email_address')::text IS NULL OR email_address ILIKE '%' || sqlc.narg('email_address')::text || '%')
    AND (sqlc.narg('last_name')::text IS NULL OR last_name ILIKE '%' || sqlc.narg('last_name')::text || '%')
ORDER BY id
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CreateMember :one
INSERT INTO library.member (
    username,
    first_name,
    last_name,
    gender,
    date_of_birth,
    member_since,
    is_student,
    email_address,
    interests
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9
)
RETURNING
    id,
    version,
    username,
    first_name,
    last_name,
    gender,
    date_of_birth,
    member_since,
    is_student,
    email_address,
    interests,
    generated,
    updated;

-- name: UpdateMember :one
UPDATE library.member
SET
    version = version + 1,
    username = $3,
    first_name = $4,
    last_name = $5,
    gender = $6,
    date_of_birth = $7,
    member_since = $8,
    is_student = $9,
    email_address = $10,
    interests = $11,
    updated = now()
WHERE id = $1
  AND version = $2
RETURNING
    id,
    version,
    username,
    first_name,
    last_name,
    gender,
    date_of_birth,
    member_since,
    is_student,
    email_address,
    interests,
    generated,
    updated;

-- name: DeleteMember :execrows
DELETE FROM library.member
WHERE id = $1;
