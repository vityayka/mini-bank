package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTPayload struct {
	Payload
	jwt.RegisteredClaims
}

func NewJWTPayload(username string, duration time.Duration) (*JWTPayload, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	issuedAt := time.Now()
	expiresAt := time.Now().Add(duration)

	payload := &JWTPayload{}

	payload.Payload = Payload{
		ID:        uuid,
		Username:  username,
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
	}

	payload.RegisteredClaims = jwt.RegisteredClaims{
		Issuer:    "mini-bank",
		Subject:   "auth",
		Audience:  jwt.ClaimStrings{username},
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		NotBefore: jwt.NewNumericDate(issuedAt),
		IssuedAt:  jwt.NewNumericDate(issuedAt),
	}

	return payload, nil
}
