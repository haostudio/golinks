package auth

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserPwd(t *testing.T) {
	u := User{}
	pwd := "hello"
	require.Error(t, u.VerifyPassword(pwd))
	require.NoError(t, u.SetPassword(pwd))
	require.NoError(t, u.VerifyPassword(pwd))
}
