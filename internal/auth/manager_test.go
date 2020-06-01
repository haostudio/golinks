package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	. "github.com/haostudio/golinks/internal/auth"
	"github.com/haostudio/golinks/internal/auth/kv"
	"github.com/haostudio/golinks/internal/encoding/gob"
	"github.com/haostudio/golinks/internal/kv/memory"
)

func TestManagerRegisterUser(t *testing.T) {
	{
		manager := testManager()
		user, err := NewUser("email@test.com", "test_pwd", "")
		require.NoError(t, err)
		require.NoError(t, manager.RegisterUser(context.Background(), *user))
	}
	{
		manager := testManager()
		user, err := NewUser("email@test.com", "test_pwd", "org")
		require.NoError(t, err)
		err = manager.RegisterUser(context.Background(), *user)
		require.True(t, errors.Is(err, ErrNotFound))
		org := Organization{
			Name:       "org",
			AdminEmail: "xxx",
		}
		require.NoError(t, manager.SetOrg(context.Background(), org))
		require.NoError(t, manager.RegisterUser(context.Background(), *user))
	}
}

func TestManagerRegisterOrg(t *testing.T) {
	{
		manager := testManager()
		org := Organization{
			Name:       "org",
			AdminEmail: "admin@test.com",
		}
		err := manager.RegisterOrg(context.Background(), org)
		require.True(t, errors.Is(err, ErrNotFound))
		user, err := NewUser("admin@test.com", "test_pwd", "")
		require.NoError(t, err)
		require.NoError(t, manager.SetUser(context.Background(), *user))
		require.NoError(t, manager.RegisterOrg(context.Background(), org))

		err = manager.RegisterOrg(context.Background(), Organization{
			Name:       "orgxx",
			AdminEmail: "admin@test.com",
		})
		require.True(t,
			errors.Is(err, ErrBadParams), "%v is not  %v", err, ErrBadParams)
	}
}

func testManager() *Manager {
	return New(Config{
		Provider:         kv.New(memory.New().In("auth"), gob.New()),
		TokenExpieration: 1 * 24 * time.Hour,
		TokenSecret:      []byte("token_secret"),
	})
}
