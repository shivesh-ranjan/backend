-- name: CreateRole :one
INSERT INTO "roles" (
	"role"
) VALUES (
	$1
) RETURNING *;

-- name: DeleteRole :exec
DELETE FROM "roles" WHERE "role"=$1;
