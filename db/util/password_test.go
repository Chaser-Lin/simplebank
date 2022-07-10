package util

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPassword(t *testing.T) {
	password := RandomString(6)

	hashedPassword1, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword1)

	err = CheckPassword(hashedPassword1, password)
	require.NoError(t, err)

	wrongPassword := RandomString(6)
	hashedWrongPassword, err := HashPassword(wrongPassword)
	require.NoError(t, err)
	require.NotEqual(t, hashedPassword1, hashedWrongPassword)

	hashedPassword2, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEqual(t, hashedPassword1, hashedPassword2)
}
