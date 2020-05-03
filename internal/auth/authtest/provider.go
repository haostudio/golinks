package authtest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/haostudio/golinks/internal/auth"
)

// ProviderLogicTest test provider logic.
func ProviderLogicTest(t *testing.T, provider auth.Provider) {
	var err error
	ctx := context.Background()

	user, err := auth.NewUser("golinks@haostudio", "pwd", "haostudio")
	require.NoError(t, err)

	_, err = provider.GetUser(ctx, user.Email)
	require.Error(t, auth.ErrNotFound, err)
	err = provider.SetUser(ctx, *user)
	require.NoError(t, err)
	u, err := provider.GetUser(ctx, user.Email)
	require.NoError(t, err)
	require.Equal(t, *user, u)

	org := auth.Organization{
		Name:       "haostudio",
		AdminEmail: "hao@haostudio",
	}
	_, err = provider.GetOrg(ctx, org.Name)
	require.Error(t, auth.ErrNotFound, err)
	err = provider.SetOrg(ctx, org)
	require.NoError(t, err)
	o, err := provider.GetOrg(ctx, org.Name)
	require.NoError(t, err)
	require.Equal(t, org, o)
}
