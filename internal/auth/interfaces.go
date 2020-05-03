package auth

import (
	"context"
	"fmt"
)

// Provider defines the auth provider interface.
type Provider interface {
	fmt.Stringer

	// users
	GetUser(ctx context.Context, email string) (User, error)
	GetUsers(ctx context.Context) ([]string, error)
	SetUser(ctx context.Context, user User) error
	DeleteUser(ctx context.Context, email string) error

	// organization
	GetOrg(ctx context.Context, name string) (Organization, error)
	GetOrgUsers(ctx context.Context, name string) ([]string, error)
	SetOrg(ctx context.Context, org Organization) error
	DeleteOrg(ctx context.Context, name string) error
}
