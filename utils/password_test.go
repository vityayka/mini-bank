package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := "password"
	hash, err := HashedPassword(password)

	require.NoError(t, err)

	// the function has to emit a different hash every call
	hash1, err := HashedPassword(password)
	require.NotEqual(t, hash1, hash)

	err = CompareHashAndPassword(hash, password)
	require.NoError(t, err)

	wrongPassword := "wrong_password"
	err = CompareHashAndPassword(hash, wrongPassword)
	require.ErrorIs(t, err, bcrypt.ErrMismatchedHashAndPassword)
}
