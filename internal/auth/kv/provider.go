package kv

import (
	"context"
	"errors"
	"fmt"

	"github.com/haostudio/golinks/internal/auth"
	"github.com/haostudio/golinks/internal/encoding"
	"github.com/haostudio/golinks/internal/kv"
)

const (
	userNamespace  = "_user"
	orgNamespace   = "_org"
	tokenNamespace = "_token"
)

// New returns an auth provider.
func New(ns kv.Namespace, enc encoding.Binary) auth.Provider {
	return &provider{
		store: ns,
		enc:   enc,
	}
}

type provider struct {
	store kv.Namespace
	enc   encoding.Binary
}

// users
func (p *provider) GetUser(ctx context.Context, email string) (
	user auth.User, err error) {
	if len(email) == 0 {
		err = fmt.Errorf("user email is required. %w", auth.ErrBadParams)
		return
	}
	// Get blob from kv
	b, err := p.store.In(userNamespace).Get(ctx, email)
	if errors.Is(err, kv.ErrNotFound) {
		err = auth.ErrNotFound
		return
	}
	if err != nil {
		err = fmt.Errorf("%v: %w", err, auth.ErrStoreError)
		return
	}
	// Decode
	err = p.enc.Decode(b, &user)
	return
}

func (p *provider) getUsersByOrg(ctx context.Context, org string) (
	users []string, err error) {
	// nolint: godox
	// FIXME: This function is too tough. Although there is an org parameter
	// support to filter result, this function still fetch all user from db and
	// filter result in code level. Maybe we could implement metadata again
	// (same as namespace of link).
	err = p.store.In(userNamespace).Iterate(ctx,
		func(key string, value []byte) bool {
			var user auth.User
			iterErr := p.enc.Decode(value, &user)
			if iterErr != nil {
				return true
			}
			if org != "" && org != user.Organization {
				return true
			}
			users = append(users, user.Email)
			return true
		})
	if errors.Is(err, kv.ErrNotFound) {
		return users, auth.ErrNotFound
	}
	return
}

func (p *provider) GetUsers(ctx context.Context) ([]string, error) {
	return p.getUsersByOrg(ctx, "")
}

func (p *provider) GetOrgUsers(ctx context.Context, name string) (
	[]string, error) {
	return p.getUsersByOrg(ctx, name)
}

func (p *provider) SetUser(ctx context.Context, user auth.User) error {
	if len(user.Email) == 0 {
		return fmt.Errorf("user email is required. %w", auth.ErrBadParams)
	}
	blob, err := p.enc.Encode(user)
	if err != nil {
		return err
	}
	return p.store.In(userNamespace).Set(ctx, user.Email, blob)
}

func (p *provider) DeleteUser(ctx context.Context, email string) error {
	return p.store.In(userNamespace).Delete(ctx, email)
}

// organization
func (p *provider) GetOrg(ctx context.Context, name string) (
	org auth.Organization, err error) {
	if len(name) == 0 {
		err = fmt.Errorf("org name is required. %w", auth.ErrBadParams)
		return
	}
	// Get blob from kv
	b, err := p.store.In(orgNamespace).Get(ctx, name)
	if errors.Is(err, kv.ErrNotFound) {
		err = auth.ErrNotFound
		return
	}
	if err != nil {
		err = fmt.Errorf("%v: %w", err, auth.ErrStoreError)
		return
	}
	// Decode
	err = p.enc.Decode(b, &org)
	return
}

func (p *provider) SetOrg(ctx context.Context, org auth.Organization) error {
	if len(org.Name) == 0 {
		return fmt.Errorf("org name is required. %w", auth.ErrBadParams)
	}
	blob, err := p.enc.Encode(org)
	if err != nil {
		return err
	}
	return p.store.In(orgNamespace).Set(ctx, org.Name, blob)
}

func (p *provider) DeleteOrg(ctx context.Context, name string) error {
	return p.store.In(orgNamespace).Delete(ctx, name)
}

// tokens
func (p *provider) GetToken(ctx context.Context, tokenStr string) (
	token auth.Token, err error) {
	if len(tokenStr) == 0 {
		err = fmt.Errorf("token is required. %w", auth.ErrBadParams)
		return
	}
	// Get blob from kv
	b, err := p.store.In(tokenNamespace).Get(ctx, tokenStr)
	if errors.Is(err, kv.ErrNotFound) {
		err = auth.ErrNotFound
		return
	}
	if err != nil {
		err = fmt.Errorf("%v: %w", err, auth.ErrStoreError)
		return
	}
	// Decode
	err = p.enc.Decode(b, &token)
	return
}

func (p *provider) SetToken(ctx context.Context, token auth.Token) error {
	if len(token.JWT) == 0 {
		return fmt.Errorf("token is required. %w", auth.ErrBadParams)
	}
	blob, err := p.enc.Encode(token)
	if err != nil {
		return err
	}
	return p.store.In(tokenNamespace).Set(ctx, token.JWT, blob)
}

func (p *provider) DeleteToken(ctx context.Context, token string) error {
	return p.store.In(tokenNamespace).Delete(ctx, token)
}

func (p *provider) String() string {
	return fmt.Sprintf("kv.provider(%s/%s)", p.store, p.enc)
}
