// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"time"
)

type Role struct {
	Role string `json:"role"`
}

type User struct {
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	Password  string    `json:"password"`
	About     string    `json:"about"`
	Photo     string    `json:"photo"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}
