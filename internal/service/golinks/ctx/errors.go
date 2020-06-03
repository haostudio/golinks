package ctx

import "errors"

// Exported errors.
var (
	ErrNotFound = errors.New("not found")
	ErrInternal = errors.New("internal")
)
