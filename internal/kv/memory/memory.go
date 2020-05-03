package memory

import (
	"fmt"

	"github.com/haostudio/golinks/internal/kv"
)

// New returns a memory store.
func New() kv.Store {
	return newStore(newMapStore())
}

// NewLRU returns a lru memory store.
func NewLRU(cap int) kv.Store {
	return newStore(newLRUStore(cap))
}

func newStore(store store) kv.Store {
	return &container{
		store: store,
	}
}

type container struct {
	store store
}

// Close finalizes a kv.Store.
func (c *container) Close() error {
	return nil
}

// In returns the namespace instance with path.
func (c *container) In(path ...string) kv.Namespace {
	return &namespace{c.store, path}
}

func (c *container) String() string {
	return fmt.Sprintf("memory.store(%s)", c.store)
}
