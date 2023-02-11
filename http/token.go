package http

import (
	"errors"
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/google/uuid"
	"github.com/o1egl/paseto"
)

var (
	ErrTokenExpired = errors.New("the PASETO has expired")
	ErrTokenInvalid = errors.New("the PASETO is not valid")
)

type PasetoPayload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issuedAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func NewPasetoPayload(username string, duration time.Duration) (*PasetoPayload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		fmt.Errorf("Could not generate a random token ID! %v", err)
		return nil, err
	}

	payload := &PasetoPayload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(duration * time.Second),
	}

	return payload, nil
}

type PasetoCreator struct {
	symmetricKey []byte
	paseto       *paseto.V2
}

func NewPasetoCreator(symmetricKey string) (*PasetoCreator, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size, must be %d bytes", chacha20poly1305.KeySize)
	}

	pc := &PasetoCreator{
		symmetricKey: []byte(symmetricKey),
		paseto:       paseto.NewV2(),
	}
	return pc, nil
}

func (c *PasetoCreator) CreateToken(username string, duration time.Duration) (string, *PasetoPayload, error) {
	payload, err := NewPasetoPayload(username, duration)
	if err != nil {
		return "", payload, err
	}

	token, err := c.paseto.Encrypt(c.symmetricKey, payload, nil)
	return token, payload, err
}

func (c *PasetoCreator) VerifyToken(token string) (*PasetoPayload, error) {
	payload := &PasetoPayload{}

	err := c.paseto.Decrypt(token, c.symmetricKey, payload, nil)
	if err != nil {
		return nil, fmt.Errorf("the PASETO is invalid. %v", ErrTokenInvalid)
	}

	if time.Now().After(payload.ExpiresAt) {
		return nil, fmt.Errorf("the PASETO has expired %v", ErrTokenExpired)
	}

	return payload, nil
}
