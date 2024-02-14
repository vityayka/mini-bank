package token

import (
	"bank/utils"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(utils.RandomString(32))
	require.NoError(t, err)

	userID := utils.RandomInt(1, 1000)
	duration := time.Minute

	issuedAt := time.Now()
	expiresAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(userID, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.NotZero(t, payload.ID)
	require.Equal(t, payload.UserID, userID)
	require.WithinDuration(t, expiresAt, payload.ExpiresAt, time.Second)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
}

func TestExpiredToken(t *testing.T) {
	maker, err := NewJWTMaker(utils.RandomString(32))
	require.NoError(t, err)

	userID := utils.RandomInt(1, 1000)
	duration := -time.Minute

	token, err := maker.CreateToken(userID, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.ErrorContains(t, err, "expired")
	require.Nil(t, payload)
}

func TestShortSecretKey(t *testing.T) {
	_, err := NewJWTMaker(utils.RandomString(31))
	require.ErrorIs(t, err, ErrSecretKeyTooShort)
}

func TestAlgNone(t *testing.T) {
	userID := utils.RandomInt(1, 1000)
	duration := time.Minute

	jwtPayload, err := NewJWTPayload(userID, duration)
	require.NoError(t, err)

	token := jwt.NewWithClaims(jwt.SigningMethodNone, jwtPayload)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	tokenString, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)

	maker, err := NewJWTMaker(utils.RandomString(32))
	require.NoError(t, err)
	_, err = maker.VerifyToken(tokenString)
	require.ErrorIs(t, err, ErrInvalidToken)
}
