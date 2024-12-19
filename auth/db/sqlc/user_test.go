package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"shivesh-ranjan.github.io/m/utils"
)

func CreateRandomUser(t *testing.T) User {
	args := CreateUserParams{
		Name:     utils.RandomString(10),
		Username: utils.RandomString(8),
		Role:     CreateRandomRole(t),
		About:    utils.RandomString(30),
		Photo:    utils.RandomString(9),
		Password: utils.RandomString(15),
	}
	user, err := testQueries.CreateUser(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, args.Role, user.Role)
	require.Equal(t, args.Name, user.Name)
	require.Equal(t, args.Username, user.Username)
	require.Equal(t, args.About, user.About)
	require.Equal(t, args.Photo, user.Photo)
	require.Equal(t, args.Password, user.Password)
	require.WithinDuration(t, user.CreatedAt, time.Now(), time.Minute)

	return user
}

func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := CreateRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)

	require.NoError(t, err)
	require.Equal(t, user1, user2)
}

func TestUpdateUser(t *testing.T) {
	user1 := CreateRandomUser(t)
	args := UpdateUserParams{
		Name:     utils.RandomString(10),
		About:    utils.RandomString(30),
		Photo:    utils.RandomString(9),
		Username: user1.Username,
	}

	user2, err := testQueries.UpdateUser(context.Background(), args)
	require.NoError(t, err)
	require.Equal(t, user2.Username, args.Username)
	require.Equal(t, user2.About, args.About)
	require.Equal(t, user2.Name, args.Name)
	require.Equal(t, user2.Photo, args.Photo)
}

func TestUpdatePassword(t *testing.T) {
	user1 := CreateRandomUser(t)
	args := UpdatePasswordParams{
		Password: utils.RandomString(12),
		Username: user1.Username,
	}

	user2, err := testQueries.UpdatePassword(context.Background(), args)
	require.NoError(t, err)
	require.Equal(t, user2.Password, args.Password)
	require.Equal(t, user2.Username, args.Username)
}

func TestUpdateRole(t *testing.T) {
	user1 := CreateRandomUser(t)
	args := UpdateRoleParams{
		Role:     CreateRandomRole(t),
		Username: user1.Username,
	}

	user2, err := testQueries.UpdateRole(context.Background(), args)
	require.NoError(t, err)
	require.Equal(t, user2.Username, args.Username)
	require.Equal(t, user2.Role, args.Role)
}
