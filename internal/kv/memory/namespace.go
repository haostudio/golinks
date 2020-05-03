package memory

import (
	"context"
	"fmt"

	"github.com/haostudio/golinks/internal/kv"
)

type namespace struct {
	store store
	path  []string
}

// In returns the namespace instance with path.
func (n *namespace) In(path ...string) kv.Namespace {
	if len(path) == 0 {
		return n
	}
	p := append(n.path[:0:0], n.path...)
	p = append(p, path...)
	return &namespace{n.store, p}
}

func (n *namespace) Get(ctx context.Context, key string) ([]byte, error) {
	keys := append(n.path[:0:0], n.path...)
	keys = append(keys, key)
	v, ok := n.store.get(keys...)
	if !ok {
		return nil, kv.ErrNotFound
	}
	return v, nil
}

func (n *namespace) Set(ctx context.Context, key string, value []byte) error {
	keys := append(n.path[:0:0], n.path...)
	keys = append(keys, key)
	n.store.set(value, keys...)
	return nil
}

func (n *namespace) Delete(ctx context.Context, key string) error {
	keys := append(n.path[:0:0], n.path...)
	keys = append(keys, key)
	n.store.del(keys...)
	return nil
}

func (n *namespace) Iterate(
	ctx context.Context, f func(key string, value []byte) (next bool)) error {
	exists, err := n.store.iter(f, n.path...)
	if err != nil {
		return err
	}
	if !exists {
		return kv.ErrNotFound
	}
	return nil
}

// Drop drops all the data in the namespace.
func (n *namespace) Drop(ctx context.Context) error {
	n.store.drop(n.path...)
	return nil
}

func (n *namespace) String() string {
	return fmt.Sprintf("memory.namespace(%s)", n.store)
}
