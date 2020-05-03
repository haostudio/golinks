package auth

import "errors"

// Exported errors.
var (
	ErrNotFound   = errors.New("not found")
	ErrStoreError = errors.New("internal store error")
	ErrBadParams  = errors.New("bad parameters")
)
