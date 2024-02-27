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

func (pm *PasetoMaker) CreateToken(userID int64, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(userID, duration)
	claims := Claims{*payload, payload.ExpiresAt}
	if err != nil {
		return "", payload, err
	}

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", payload, err
	}

	token, err := paseto.NewTokenFromClaimsJSON(claimsJSON, []byte(""))
	if err != nil {
		return "", payload, err
	}

	return token.V4Encrypt(pm.symmetricKey, []byte(implicitString)), payload, nil
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
