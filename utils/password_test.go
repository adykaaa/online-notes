package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashUserPassword(t *testing.T) {
	t.Run("hashes password OK", func(t *testing.T) {
		hpw, err := HashUserPassword("abc123!")
		require.NoError(t, err)
		require.NotEmpty(t, hpw)
	})
	t.Run("fails if too short", func(t *testing.T) {
		hpw, err := HashUserPassword("abc1")
		require.Error(t, err)
		require.Empty(t, hpw)
	})
	t.Run("fails if empty", func(t *testing.T) {
		hpw, err := HashUserPassword("")
		require.Error(t, err)
		require.Empty(t, hpw)
	})
	t.Run("fails if too long", func(t *testing.T) {
		hpw, err := HashUserPassword(".vI(5dSO^hM)Q:>z.n'T?1mdzFQE2;UP5N-(q`NCkM=m'efZZ'JajBn№A)vU:84Mozt<G:vg*")
		require.Error(t, err)
		require.Empty(t, hpw)
	})
}

func TestValidateUserPassword(t *testing.T) {

}
