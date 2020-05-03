package cached

import (
	"context"
	"fmt"

	"github.com/haostudio/golinks/internal/kv"
)

type namespace struct {
	store     *store
	root      []string
	canonical kv.Namespace
	cache     kv.Namespace
}

// In returns the namespace instance with path.
func (n *namespace) In(path ...string) kv.Namespace {
	p := append(n.root[:0:0], n.root...)
	p = append(p, path...)
	return n.store.In(p...)
}

// Get returns the value in the namespace with key.
func (n *namespace) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := n.cache.Get(ctx, key)
	if err == nil {
		return val, nil
	}
	val, err = n.canonical.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	// best effort to set the cache
	cloned := append(val[:0:0], val...)
	_ = n.cache.Set(ctx, key, cloned)
	return val, nil
}

// Set sets the value in the namespace with key.
func (n *namespace) Set(ctx context.Context, key string, value []byte) error {
	err := n.canonical.Set(ctx, key, value)
	if err != nil {
		return err
	}
	// best effort to set the cache
	cloned := append(value[:0:0], value...)
	_ = n.cache.Set(ctx, key, cloned)
	return nil
}

// Delete deletes the value in the namespace with key.
func (n *namespace) Delete(ctx context.Context, key string) error {
	err := n.canonical.Delete(ctx, key)
	if err != nil {
		return err
	}
	// best effort to delete cache
	_ = n.cache.Delete(ctx, key)
	return nil
}

// Iterate iterates the values in the namespace.
func (n *namespace) Iterate(
	ctx context.Context, f func(key string, value []byte) (next bool)) error {
	err := n.cache.Iterate(ctx, f)
	if err == nil {
		return nil
	}
	return n.canonical.Iterate(ctx, f)
}

// Drop drops all the values in the namespace.
func (n *namespace) Drop(ctx context.Context) error {
	err := n.canonical.Drop(ctx)
	if err != nil {
		return err
	}
	// best effort to delete cache
	_ = n.cache.Drop(ctx)
	return nil
}

func (n *namespace) String() string {
	return fmt.Sprintf("cached.namespace(%s/%s)", n.canonical, n.cache)
}
