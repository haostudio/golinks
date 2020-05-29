package auth

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Config defines the auth manager config.
type Config struct {
	Provider         Provider
	TokenExpieration time.Duration
	TokenSecret      []byte
}

// New returns an auth manager with provider.
func New(config Config) *Manager {
	return &Manager{
		Provider:         config.Provider,
		TokenExpieration: config.TokenExpieration,
		TokenSecret:      config.TokenSecret,
	}
}

// Manager manages the authentication.
type Manager struct {
	Provider
	TokenExpieration time.Duration
	TokenSecret      []byte
}

// RegisterUser creates the user, ensuring the user does not exist and the org
// does exist if user.Organization is not empty.
func (m *Manager) RegisterUser(ctx context.Context, user User) (err error) {
	var exists bool
	exists, err = m.IsUserExists(ctx, user.Email)
	if err != nil {
		err = fmt.Errorf("failed to get user exists. %w", err)
		return
	}
	if exists {
		err = ErrUserExists
		return
	}
	if user.Organization != "" {
		exists, err = m.IsOrgExists(ctx, user.Organization)
		if err != nil {
			err = fmt.Errorf("failed to get org exists. %w", err)
			return
		}
		if !exists {
			err = fmt.Errorf("org not found. %w", ErrNotFound)
			return
		}
	}
	return m.SetUser(ctx, user)
}

// RegisterOrg creates the org.
func (m *Manager) RegisterOrg(ctx context.Context, org Organization) (
	err error) {
	var exists bool
	exists, err = m.IsOrgExists(ctx, org.Name)
	if err != nil {
		err = fmt.Errorf("failed to get org exists. %w", err)
		return
	}
	if exists {
		err = ErrOrgExists
		return
	}
	if org.AdminEmail != "" {
		exists, err = m.IsUserExists(ctx, org.AdminEmail)
		if err != nil {
			err = fmt.Errorf("failed to get admin exists. %w", err)
			return
		}
		if !exists {
			err = fmt.Errorf("admin user not found. %w", ErrNotFound)
			return
		}
	}
	// create org.
	err = m.SetOrg(ctx, org)
	if err != nil {
		err = fmt.Errorf("%v;%w", err, ErrStoreError)
	}
	return
}

// RegisterOrgWithAdmin creates the org and the admin user.
func (m *Manager) RegisterOrgWithAdmin(
	ctx context.Context, org Organization, admin User) (err error) {
	// check parameters
	if org.AdminEmail != admin.Email {
		err = ErrBadParams
		return
	}
	// check if org or user exists
	var exists bool
	exists, err = m.IsOrgExists(ctx, org.Name)
	if err != nil {
		err = fmt.Errorf("failed to get org exists. %w", err)
		return
	}
	if exists {
		err = ErrOrgExists
		return
	}
	exists, err = m.IsUserExists(ctx, admin.Email)
	if err != nil {
		err = fmt.Errorf("failed to get admin exists. %w", err)
		return
	}
	if exists {
		err = ErrUserExists
		return
	}
	// create org and user.
	err = m.SetOrg(ctx, org)
	if err != nil {
		err = fmt.Errorf("%v;%w", err, ErrStoreError)
		return
	}
	err = m.SetUser(ctx, admin)
	if err != nil {
		err = fmt.Errorf("%v;%w", err, ErrStoreError)
		// best effort to delete org
		_ = m.DeleteOrg(ctx, org.Name)
	}
	return
}

// SetUserOrg sets the org of user with email.
func (m *Manager) SetUserOrg(ctx context.Context, email, org string) error {
	user, err := m.GetUser(ctx, email)
	if err != nil {
		return fmt.Errorf("user not found. %w", err)
	}
	if user.Organization == org {
		return nil
	}
	_, err = m.GetOrg(ctx, org)
	if err != nil {
		return fmt.Errorf("org not found. %w", err)
	}
	user.Organization = org
	return m.SetUser(ctx, user)
}

// IsUserExists returns if the user exists.
func (m *Manager) IsUserExists(ctx context.Context, email string) (
	exists bool, err error) {
	_, err = m.GetUser(ctx, email)
	if err == nil {
		exists = true
		return
	}
	if errors.Is(err, ErrNotFound) {
		exists = false
		err = nil
		return
	}
	return
}

// IsOrgExists returns if the org exists.
func (m *Manager) IsOrgExists(ctx context.Context, org string) (
	exists bool, err error) {
	_, err = m.GetOrg(ctx, org)
	if err == nil {
		exists = true
		return
	}
	if errors.Is(err, ErrNotFound) {
		exists = false
		err = nil
		return
	}
	return
}

// Login verify the user's email, password and returns the access token.
func (m *Manager) Login(ctx context.Context, email string, password string) (
	token *Token, err error) {
	var user User
	user, err = m.Provider.GetUser(ctx, email)
	if err != nil {
		return
	}
	err = user.VerifyPassword(password)
	if err != nil {
		err = fmt.Errorf("%v; %w", err, ErrBadParams)
		return
	}
	token, err = NewToken(user, m.TokenSecret, m.TokenExpieration)
	if err != nil {
		return
	}
	err = m.SetToken(ctx, *token)
	if err != nil {
		token = nil
	}
	return
}

// Verify verifies the access token.
func (m *Manager) Verify(ctx context.Context, tokenStr string) (
	claims *TokenClaims, err error) {
	var token Token
	token, err = m.GetToken(ctx, tokenStr)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			err = ErrInvalidToken
			return
		}
		return
	}
	claims, err = verifyToken(token.JWT, m.TokenSecret)
	if err != nil {
		return
	}
	if claims.ExpiredAt < time.Now().Unix() {
		err = ErrTokenExpired
	}
	return
}

// Logout deletes the access token.
func (m *Manager) Logout(ctx context.Context, token string) (err error) {
	return m.DeleteToken(ctx, token)
}
