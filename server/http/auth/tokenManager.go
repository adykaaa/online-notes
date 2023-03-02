package auth

import "time"

type TokenManager interface {
	CreateToken(username string, duration time.Duration) (string, *PasetoPayload, error)
	VerifyToken(token string) (*PasetoPayload, error)
}
