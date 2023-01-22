package http

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/google/uuid"
	"github.com/o1egl/paseto"
)

type PasetoPayload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issuedAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func NewPasetoPayload(username string) (*PasetoPayload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		fmt.Errorf("Could not generate a random token ID! %v", err)
		return nil, err
	}

	payload := &PasetoPayload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	return payload, nil
}

type Paseto struct {
	key    []byte
	paseto *paseto.V2
}

func NewPaseto(symmetricKey string) (*Paseto, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size, must be %d bytes", chacha20poly1305.KeySize)
	}

	paseto := &Paseto{
		key:    []byte(symmetricKey),
		paseto: paseto.NewV2(),
	}
	return paseto, nil
}

func (maker *Paseto) CreateToken(username string) (string, *PasetoPayload, error) {
	payload, err := NewPasetoPayload(username)
	if err != nil {
		return "", payload, err
	}

	token, err := maker.paseto.Encrypt(maker.key, payload, nil)
	return token, payload, err
}

func (maker *Paseto) VerifyToken(token string) (*PasetoPayload, error) {
	payload := &PasetoPayload{}

	err := maker.paseto.Decrypt(token, maker.key, payload, nil)
	if err != nil {
		return nil, fmt.Errorf("the PASETO is invalid!")
	}

	if time.Now().After(payload.ExpiresAt) {
		return nil, fmt.Errorf("the PASETO has expired")
	}

	return payload, nil
}
