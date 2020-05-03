package kv

import "errors"

// Exported errors
var (
	ErrNotFound      = errors.New("value not found")
	ErrInternalError = errors.New("internal store error")
	ErrNotSupport    = errors.New("not support")
)
