package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"shivesh-ranjan.github.io/m/utils"
)

func CreateRandomRole(t *testing.T) string {
	roleName := utils.RandomString(5)
	role, err := testQueries.CreateRole(context.Background(), roleName)

	require.NoError(t, err)
	require.Equal(t, roleName, role)

	return role
}

func TestCreateRole(t *testing.T) {
	CreateRandomRole(t)
}

func TestDeleteRole(t *testing.T) {
	role := CreateRandomRole(t)

	testQueries.DeleteRole(context.Background(), role)

	user := CreateRandomUser(t)
	args := UpdateRoleParams{
		Role:     role,
		Username: user.Username,
	}
	user1, err := testQueries.UpdateRole(context.Background(), args)
	require.Error(t, err)
	require.Empty(t, user1)
}
