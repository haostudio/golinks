package link

import "errors"

// Exported errors.
var (
	ErrVersionNotSupport = errors.New("link version is not support")
	ErrInvalidParams     = errors.New("invalid parameters passed")
	ErrNotFound          = errors.New("link not found")
)
