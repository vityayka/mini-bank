package token

import "time"

const minSecretKeySize = 32

type Maker interface {
	CreateToken(userID int64, duration time.Duration) (string, error)

	VerifyToken(token string) (*Payload, error)
}
