// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"context"
)

type Querier interface {
	CreateRole(ctx context.Context, role string) (string, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteRole(ctx context.Context, role string) error
	GetUser(ctx context.Context, username string) (User, error)
	UpdatePassword(ctx context.Context, arg UpdatePasswordParams) (User, error)
	UpdateRole(ctx context.Context, arg UpdateRoleParams) (User, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
}

var _ Querier = (*Queries)(nil)
