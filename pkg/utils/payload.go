package utils

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

type Payload struct {
	Id        uuid.UUID `json:"id"`
	UserId    int64     `json:"user_id"`
	Email     string    `json:"email"`
	UserType  string    `json:"user_type"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(params *TokenParams) (*Payload, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	payload := &Payload{
		Id:        tokenId,
		UserId:    params.UserId,
		Email:     params.Email,
		UserType:  params.UserType,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(params.Duration),
	}
	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
