// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: user.sql

package db

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
	username,
	name,
	password,
	about,
	photo,
	role
) VALUES (
	$1, $2, $3, $4, $5, $6
) RETURNING username, name, password, about, photo, role, created_at
`

type CreateUserParams struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Password string `json:"password"`
	About    string `json:"about"`
	Photo    string `json:"photo"`
	Role     string `json:"role"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Username,
		arg.Name,
		arg.Password,
		arg.About,
		arg.Photo,
		arg.Role,
	)
	var i User
	err := row.Scan(
		&i.Username,
		&i.Name,
		&i.Password,
		&i.About,
		&i.Photo,
		&i.Role,
		&i.CreatedAt,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT username, name, password, about, photo, role, created_at FROM users
WHERE username=$1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, username)
	var i User
	err := row.Scan(
		&i.Username,
		&i.Name,
		&i.Password,
		&i.About,
		&i.Photo,
		&i.Role,
		&i.CreatedAt,
	)
	return i, err
}

const updatePassword = `-- name: UpdatePassword :one
UPDATE users SET password=$1 WHERE username=$2 RETURNING username, name, password, about, photo, role, created_at
`

type UpdatePasswordParams struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

func (q *Queries) UpdatePassword(ctx context.Context, arg UpdatePasswordParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updatePassword, arg.Password, arg.Username)
	var i User
	err := row.Scan(
		&i.Username,
		&i.Name,
		&i.Password,
		&i.About,
		&i.Photo,
		&i.Role,
		&i.CreatedAt,
	)
	return i, err
}

const updateRole = `-- name: UpdateRole :one
UPDATE users SET role=$1 WHERE username=$2 RETURNING username, name, password, about, photo, role, created_at
`

type UpdateRoleParams struct {
	Role     string `json:"role"`
	Username string `json:"username"`
}

func (q *Queries) UpdateRole(ctx context.Context, arg UpdateRoleParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateRole, arg.Role, arg.Username)
	var i User
	err := row.Scan(
		&i.Username,
		&i.Name,
		&i.Password,
		&i.About,
		&i.Photo,
		&i.Role,
		&i.CreatedAt,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :one
UPDATE users SET name=$1, about=$2, photo=$3 WHERE username=$4 RETURNING username, name, password, about, photo, role, created_at
`

type UpdateUserParams struct {
	Name     string `json:"name"`
	About    string `json:"about"`
	Photo    string `json:"photo"`
	Username string `json:"username"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUser,
		arg.Name,
		arg.About,
		arg.Photo,
		arg.Username,
	)
	var i User
	err := row.Scan(
		&i.Username,
		&i.Name,
		&i.Password,
		&i.About,
		&i.Photo,
		&i.Role,
		&i.CreatedAt,
	)
	return i, err
}
