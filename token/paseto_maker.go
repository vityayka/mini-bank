package token

import (
	"encoding/json"
	"fmt"
	"time"

	"aidanwoods.dev/go-paseto"
)

const implicitString = "azazaz nahooy lalka"

type PasetoMaker struct {
	symmetricKey paseto.V4SymmetricKey
}

type Claims struct {
	Payload
	Exp time.Time `json:"exp"`
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) < 32 {
		return nil, fmt.Errorf("symmetric key size should be of %d bytes", minSecretKeySize)
	}
	v4SymmetricKey, err := paseto.V4SymmetricKeyFromBytes([]byte(symmetricKey))
	if err != nil {
		return nil, err
	}
	return &PasetoMaker{v4SymmetricKey}, nil
}

func (pm *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	claims := Claims{*payload, payload.ExpiresAt}
	if err != nil {
		return "", err
	}

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	token, err := paseto.NewTokenFromClaimsJSON(claimsJSON, []byte(""))
	if err != nil {
		return "", err
	}

	return token.V4Encrypt(pm.symmetricKey, []byte(implicitString)), nil
}

func (pm *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	pToken, err := paseto.NewParser().ParseV4Local(pm.symmetricKey, token, []byte(implicitString))
	if err != nil {
		return nil, err
	}

	claims := Claims{}

	err = json.Unmarshal(pToken.ClaimsJSON(), &claims)
	return &claims.Payload, err
}
