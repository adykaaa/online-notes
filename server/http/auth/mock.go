package auth

import (
	"time"
)

type MockTokenManager struct {
	ReturnInvalidToken bool
	ReturnExpiredToken bool
}

func (m *MockTokenManager) CreateToken(username string, duration time.Duration) (string, *PasetoPayload, error) {
	return "testtoken",
		&PasetoPayload{},
		nil
}

func (m *MockTokenManager) VerifyToken(token string) (*PasetoPayload, error) {
	if m.ReturnExpiredToken {
		return nil, ErrTokenExpired
	}
	if m.ReturnInvalidToken {
		return nil, ErrTokenInvalid
	}
	return &PasetoPayload{},
		nil
}
