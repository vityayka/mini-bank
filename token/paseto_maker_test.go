package token

import (
	"bank/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(utils.RandomString(32))
	require.NoError(t, err)

	userID := utils.RandomInt(1, 1000)
	duration := time.Minute

	issuedAt := time.Now()
	expiresAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(userID, utils.Depositor, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.NotZero(t, payload.ID)
	require.Equal(t, payload.UserID, userID)
	require.WithinDuration(t, expiresAt, payload.ExpiresAt, time.Second)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(utils.RandomString(32))
	require.NoError(t, err)

	userID := utils.RandomInt(1, 1000)
	duration := -time.Minute

	token, payload, err := maker.CreateToken(userID, utils.Depositor, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err = maker.VerifyToken(token)
	require.ErrorContains(t, err, "expired")
	require.Nil(t, payload)
}

func TestPasetoShortSecretKey(t *testing.T) {
	_, err := NewJWTMaker(utils.RandomString(31))
	require.ErrorIs(t, err, ErrSecretKeyTooShort)
}
