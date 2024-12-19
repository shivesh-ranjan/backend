-- name: CreateUser :one
INSERT INTO users (
	username,
	name,
	password,
	about,
	photo,
	role
) VALUES (
	$1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username=$1 LIMIT 1;

-- name: UpdateUser :one
UPDATE users SET name=$1, about=$2, photo=$3 WHERE username=$4 RETURNING *;

-- name: UpdatePassword :one
UPDATE users SET password=$1 WHERE username=$2 RETURNING *;

-- name: UpdateRole :one
UPDATE users SET role=$1 WHERE username=$2 RETURNING *;
