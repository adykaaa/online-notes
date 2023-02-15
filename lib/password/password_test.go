package password

import (
	"testing"

	"github.com/adykaaa/online-notes/lib/random"
	"github.com/stretchr/testify/require"
)

func TestHashUserPassword(t *testing.T) {
	t.Run("password hashing OK", func(t *testing.T) {
		hpw, err := Hash(random.NewString(10))
		require.NoError(t, err)
		require.NotEmpty(t, hpw)
	})

	t.Run("same hash for different password", func(t *testing.T) {
		pw := random.NewString(20)
		hpw1, err := Hash(pw)
		require.NoError(t, err)
		hpw2, err := Hash(pw)
		require.NoError(t, err)
		require.NotEqual(t, hpw1, hpw2)
	})
	t.Run("different hash for different password", func(t *testing.T) {
		pw := random.NewString(20)
		hpw1, err := Hash(pw)
		require.NoError(t, err)
		hpw2, err := Hash(pw)
		require.NoError(t, err)
		require.NotEqual(t, hpw1, hpw2)
	})
	t.Run("fails if too short", func(t *testing.T) {
		hpw, err := Hash(random.NewString(4))
		require.Error(t, err)
		require.Empty(t, hpw)
	})
	t.Run("fails if empty", func(t *testing.T) {
		hpw, err := Hash("")
		require.Error(t, err)
		require.Empty(t, hpw)
	})
	t.Run("fails if too long", func(t *testing.T) {
		hpw, err := Hash(random.NewString(100))
		require.Error(t, err)
		require.Empty(t, hpw)
	})

}

func TestValidate(t *testing.T) {
	const pw1 = "abc123!"
	const pw2 = "abc321!"
	t.Run("password validation OK", func(t *testing.T) {
		hpw, _ := Hash(pw1)
		err := Validate(hpw, pw1)
		require.NoError(t, err)
	})
	t.Run("fails with not same hash", func(t *testing.T) {
		hpw, _ := Hash(pw1)
		err := Validate(hpw, pw2)
		require.Error(t, err)
	})
}
