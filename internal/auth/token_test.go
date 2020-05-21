package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestToken(t *testing.T) {
	secret := []byte("test secret")
	now := time.Now()
	params := tokenParams{
		User: User{
			Email:        "hi",
			PasswordHash: []byte("how are you"),
			Organization: "hello",
		},
		IssuedAt:  now,
		ExpiredAt: now.Add(1),
		Secret:    secret,
	}
	token, err := genToken(params)
	require.NoError(t, err)
	claims, err := verifyToken(token, secret)
	require.NoError(t, err)
	require.Equal(t, "hi", claims.Email)
	require.Equal(t, "hello", claims.Org)
	require.Equal(t, now.Unix(), claims.IssuedAt)
	require.Equal(t, now.Add(1).Unix(), claims.ExpiredAt)
}
