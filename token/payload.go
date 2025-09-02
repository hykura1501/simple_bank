package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

var ErrExpiredToken = errors.New("token has expired")
var ErrInvalidToken = errors.New("token is invalid")

type Payload struct {
	ID        pgtype.UUID `json:"id"`
	Username  string      `json:"username"`
	IssuedAt  time.Time   `json:"issued_at"`
	ExpiredAt time.Time   `json:"expired_at"`
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	u := uuid.New()
	pgUUID := pgtype.UUID{
		Bytes: u,
		Valid: true,
	}
	payload := &Payload{
		ID:        pgUUID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}

func (p *Payload) Valid() error {
	if p.ExpiredAt.Before(time.Now()) {
		return ErrExpiredToken
	}
	return nil
}
