package token

import (
	"bank/utils"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTMaker struct {
	secretKey string
}

var (
	ErrInvalidToken      = errors.New("invalid token")
	ErrSecretKeyTooShort = errors.New("secret key too short")
)

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, ErrSecretKeyTooShort
	}

	return &JWTMaker{secretKey}, nil
}

type JWTPayload struct {
	Payload
	jwt.RegisteredClaims
}

func NewJWTPayload(userID int64, role utils.Role, duration time.Duration) (*JWTPayload, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	issuedAt := time.Now()
	expiresAt := time.Now().Add(duration)

	payload := &JWTPayload{}

	payload.Payload = Payload{
		ID:        uuid,
		UserID:    userID,
		Role:      role,
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
	}

	payload.RegisteredClaims = jwt.RegisteredClaims{
		Issuer:    "mini-bank",
		Subject:   "auth",
		Audience:  jwt.ClaimStrings{fmt.Sprintf("%d", userID)},
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		NotBefore: jwt.NewNumericDate(issuedAt),
		IssuedAt:  jwt.NewNumericDate(issuedAt),
	}

	return payload, nil
}

func (maker *JWTMaker) CreateToken(userID int64, role utils.Role, duration time.Duration) (string, *Payload, error) {
	payload, err := NewJWTPayload(userID, role, duration)
	if err != nil {
		return "", &payload.Payload, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS512, payload)
	tokenString, err := jwtToken.SignedString([]byte(maker.secretKey))

	return tokenString, &payload.Payload, err
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, isOk := token.Method.(*jwt.SigningMethodHMAC); !isOk {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &JWTPayload{}, keyFunc)
	if err != nil {
		return nil, err
	}

	payload, isOk := jwtToken.Claims.(*JWTPayload)
	if !isOk {
		return nil, ErrInvalidToken
	}

	return &payload.Payload, nil
}
