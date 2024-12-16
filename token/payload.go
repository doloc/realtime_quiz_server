package token

import (
	"errors"
	"time"
)

var (
	ErrExpiredToken = errors.New("token has expired")
	ErrInvalidToken = errors.New("token is invalid")
)

type Payload struct {
	Username  string `json:"username"`
	IssuedAt  int64  `json:"issued_at"`
	ExpiredAt int64  `json:"expired_at"`
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	payload := &Payload{
		Username:  username,
		IssuedAt:  time.Now().Unix(),
		ExpiredAt: time.Now().Add(duration).Unix(),
	}
	return payload, nil
}
