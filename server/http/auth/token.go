package auth

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
	ErrTokenMissing = errors.New("the PASETO is missing")
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

type PasetoManager struct {
	symmetricKey []byte
	paseto       *paseto.V2
}

func NewPasetoManager(symmetricKey string) (*PasetoManager, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size, must be %d bytes", chacha20poly1305.KeySize)
	}

	pc := &PasetoManager{
		symmetricKey: []byte(symmetricKey),
		paseto:       paseto.NewV2(),
	}
	return pc, nil
}

func (c *PasetoManager) CreateToken(username string, duration time.Duration) (string, *PasetoPayload, error) {
	payload, err := NewPasetoPayload(username, duration)
	if err != nil {
		return "", payload, err
	}

	token, err := c.paseto.Encrypt(c.symmetricKey, payload, nil)
	return token, payload, err
}

func (c *PasetoManager) VerifyToken(token string) (*PasetoPayload, error) {
	if token == "" {
		return nil, ErrTokenMissing
	}

	payload := &PasetoPayload{}

	err := c.paseto.Decrypt(token, c.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrTokenInvalid
	}

	if time.Now().After(payload.ExpiresAt) {
		return nil, ErrTokenExpired
	}

	return payload, nil
}
