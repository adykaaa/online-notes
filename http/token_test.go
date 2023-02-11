package http

import (
	"testing"
	"time"

	"github.com/adykaaa/online-notes/utils"
	"github.com/stretchr/testify/require"
)

func TestPasetoCreator(t *testing.T) {

	key := utils.NewRandomString(32)
	uname := utils.NewRandomString(15)
	duration := 10 * time.Second

	pc, err := NewPasetoCreator(key)
	require.NoError(t, err)

	token, payload, err := pc.CreateToken(uname, duration)
	require.NoError(t, err)
	require.Equal(t, payload.Username, uname)
	require.NotEmpty(t, token)

}
