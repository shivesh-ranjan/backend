// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: role.sql

package db

import (
	"context"
)

const createRole = `-- name: CreateRole :one
INSERT INTO "roles" (
	"role"
) VALUES (
	$1
) RETURNING role
`

func (q *Queries) CreateRole(ctx context.Context, role string) (string, error) {
	row := q.db.QueryRowContext(ctx, createRole, role)
	err := row.Scan(&role)
	return role, err
}

const deleteRole = `-- name: DeleteRole :exec
DELETE FROM "roles" WHERE "role"=$1
`

func (q *Queries) DeleteRole(ctx context.Context, role string) error {
	_, err := q.db.ExecContext(ctx, deleteRole, role)
	return err
}
