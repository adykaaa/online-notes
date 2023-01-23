package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashUserPassword(t *testing.T) {
	t.Run("should hash the password without error", func(t *testing.T) {
		pw := "abc123!"
		_, err := HashUserPassword(pw)
		require.NoError(t, err)
	})
	t.Run("should return error - too short password", func(t *testing.T) {
		pw := "abc1"
		_, err := HashUserPassword(pw)
		require.Error(t, err)
	})
	t.Run("should return error - empty password", func(t *testing.T) {
		pw := ""
		_, err := HashUserPassword(pw)
		require.Error(t, err)
	})
	t.Run("should return error - too long password", func(t *testing.T) {
		pw := ".vI(5dSO^hM)Q:>z.n'T?1mdzFQE2;UP5N-(q`NCkM=m'efZZ'JajBn№A)vU:84Mozt<G:vg*"
		_, err := HashUserPassword(pw)
		require.Error(t, err)
	})
}

func TestValidateUserPassword(t *testing.T) {

}
