package auth

import (
	"testing"
	"time"

	"github.com/adykaaa/online-notes/lib/random"
	"github.com/stretchr/testify/require"
)

func TestPasetoManager(t *testing.T) {

	key := random.NewString(32)
	uname := random.NewString(15)
	duration := 1000 * time.Second
	pc, err := NewPasetoManager(key)

	require.NoError(t, err)

	t.Run("tokenCreation and verification OK", func(t *testing.T) {
		token, payload, err := pc.CreateToken(uname, duration)

		require.NoError(t, err)
		require.Equal(t, payload.Username, uname)
		require.NotEmpty(t, token)

		payload, err = pc.VerifyToken(token)
		require.NoError(t, err)
		require.Equal(t, payload.Username, uname)
	})

	t.Run("fails because of invalid key length", func(t *testing.T) {
		pc, err := NewPasetoManager("wrongkeylength")
		require.ErrorIs(t, err, ErrInvalidSymmetricKeySize)
		require.Nil(t, pc)
	})

	t.Run("fails with invalid token", func(t *testing.T) {
		payload, err := pc.VerifyToken("invalidtoken")
		require.ErrorIs(t, err, ErrTokenInvalid)
		require.Nil(t, payload)
	})

	t.Run("fails with missing token", func(t *testing.T) {
		payload, err := pc.VerifyToken("")
		require.ErrorIs(t, err, ErrTokenMissing)
		require.Nil(t, payload)
	})

	t.Run("fails with expired token", func(t *testing.T) {
		token, _, err := pc.CreateToken(uname, 0)
		require.NoError(t, err)
		require.NotEmpty(t, token)
		retToken, err := pc.VerifyToken(token)
		require.ErrorIs(t, err, ErrTokenExpired)
		require.Nil(t, retToken)
	})

}
