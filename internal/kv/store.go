package kv

import (
	"context"
	"fmt"
)

// Store defines a key-value store interface.
type Store interface {
	fmt.Stringer

	// Close finalizes a kv.Store.
	Close() error
	// In returns the namespace instance with path.
	In(path ...string) Namespace
}

// Namespace defines a namespace in a kv.Store.
type Namespace interface {
	fmt.Stringer

	// In returns the namespace instance with path.
	In(path ...string) Namespace
	// Get returns the value in the namespace with key.
	Get(ctx context.Context, key string) ([]byte, error)
	// Set sets the value in the namespace with key.
	Set(ctx context.Context, key string, value []byte) error
	// Delete deletes the value in the namespace with key.
	Delete(ctx context.Context, key string) error
	// Iterate iterates the values in the namespace.
	Iterate(
		ctx context.Context, f func(key string, value []byte) (next bool)) error
	// Drop drops all the values in the namespace.
	Drop(ctx context.Context) error
}
