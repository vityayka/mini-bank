package token

import (
	"bank/utils"
	"time"
)

const minSecretKeySize = 32

type Maker interface {
	CreateToken(userID int64, role utils.Role, duration time.Duration) (string, *Payload, error)

	VerifyToken(token string) (*Payload, error)
}
