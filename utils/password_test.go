package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashUserPassword(t *testing.T) {
	password := "abc123"
	hashedPassword, err := HashUserPassword(password)

	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)
}

func TestValidateUserPassword(t *testing.T) {
	tc := []struct{
		name string
		password string
		hashedPassword string
	}{
		{
			name: "Validation is OK",
			password: "abc123",
			hashedPassword: HashUserPassword(password),
		}
	}
}
