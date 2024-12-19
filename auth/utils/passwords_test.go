package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPassword(t *testing.T) {
	password := RandomString(6)

	result, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	error := CheckPassword(password, result)
	require.NoError(t, error)

	wrong := RandomString(6)
	err = CheckPassword(wrong, result)
	require.Error(t, err)
}
