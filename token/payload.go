package token

import (
	"bank/utils"
	"time"

	"github.com/google/uuid"
)

type Payload struct {
	ID        uuid.UUID  `json:"id"`
	UserID    int64      `json:"user_id"`
	Role      utils.Role `json:"role"`
	IssuedAt  time.Time  `json:"issued_at"`
	ExpiresAt time.Time  `json:"expires_at"`
}

func NewPayload(userID int64, role utils.Role, duration time.Duration) (*Payload, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &Payload{
		ID:        uuid,
		UserID:    userID,
		Role:      role,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}, nil
}
