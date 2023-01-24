package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashUserPassword(t *testing.T) {
	t.Run("password hashing OK", func(t *testing.T) {
		hpw, err := HashPassword("abc123!")
		require.NoError(t, err)
		require.NotEmpty(t, hpw)
	})

	t.Run("different hash for different passwords", func(t *testing.T) {
		hpw1, err := HashPassword("abc123!")
		require.NoError(t, err)
		hpw2, err := HashPassword("abc321!")
		require.NoError(t, err)
		require.NotEqual(t, hpw1, hpw2)
	})
	t.Run("fails if too short", func(t *testing.T) {
		hpw, err := HashPassword("abc1")
		require.Error(t, err)
		require.Empty(t, hpw)
	})
	t.Run("fails if empty", func(t *testing.T) {
		hpw, err := HashPassword("")
		require.Error(t, err)
		require.Empty(t, hpw)
	})
	t.Run("fails if too long", func(t *testing.T) {
		hpw, err := HashPassword(".vI(5dSO^hM)Q:>z.n'T?1mdzFQE2;UP5N-(q`NCkM=m'efZZ'JajBn№A)vU:84Mozt<G:vg*")
		require.Error(t, err)
		require.Empty(t, hpw)
	})

}

func TestValidatePassword(t *testing.T) {
	const pw1 = "abc123!"
	const pw2 = "abc321!"
	t.Run("password validation OK", func(t *testing.T) {
		hpw, _ := HashPassword(pw1)
		err := ValidatePassword(hpw, pw1)
		require.NoError(t, err)
	})
	t.Run("fails with not same hash", func(t *testing.T) {
		hpw, _ := HashPassword(pw1)
		err := ValidatePassword(hpw, pw2)
		require.Error(t, err)
	})
}
