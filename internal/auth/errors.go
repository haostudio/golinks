package auth

import "errors"

// Exported errors.
var (
	ErrBadParams     = errors.New("bad parameters")
	ErrInternalError = errors.New("internal error")
	ErrNotFound      = errors.New("not found")

	ErrStoreError = errors.New("internal store error")
	ErrOrgExists  = errors.New("org already exists")
	ErrUserExists = errors.New("user already token")

	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token expired")
)
