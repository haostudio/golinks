package leveldb

import (
	"context"
	"errors"
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"

	"github.com/haostudio/golinks/internal/kv"
)

type namespace struct {
	store     *store
	namespace []string
}

// In returns the namespace instance with path.
func (n *namespace) In(path ...string) kv.Namespace {
	if len(path) == 0 {
		return n
	}
	ns := append(n.namespace[:0:0], n.namespace...)
	ns = append(ns, path...)
	return &namespace{
		store:     n.store,
		namespace: ns,
	}
}

// Get returns the value in the namespace with key.
func (n *namespace) Get(ctx context.Context, key string) ([]byte, error) {
	k := n.store.meta.getKeyIn(key, n.namespace...)
	var b []byte
	err := n.store.read(func(db *leveldb.DB) error {
		var readErr error
		b, readErr = n.get(db, k)
		return readErr
	})
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (n *namespace) get(db *leveldb.DB, key []byte) ([]byte, error) {
	b, err := db.Get(key, nil)
	if err != nil && errors.Is(err, leveldb.ErrNotFound) {
		return nil, fmt.Errorf("%v: %w", err, kv.ErrNotFound)
	} else if err != nil {
		return nil, fmt.Errorf("%v: %w", err, kv.ErrInternalError)
	}
	return b, nil
}

// Set sets the value in the namespace with key.
func (n *namespace) Set(ctx context.Context, key string, value []byte) error {
	return n.store.write(func(tx *leveldb.Transaction) error {
		k, err := n.store.meta.addKeyIn(tx, key, n.namespace...)
		if err != nil {
			return fmt.Errorf("%v: %w", err, kv.ErrInternalError)
		}
		err = tx.Put(k, value, &opt.WriteOptions{Sync: true})
		if err != nil {
			return fmt.Errorf("%v: %w", err, kv.ErrInternalError)
		}
		return nil
	})
}

// Delete deletes the value in the namespace with key.
func (n *namespace) Delete(ctx context.Context, key string) error {
	return n.store.write(func(tx *leveldb.Transaction) error {
		k, err := n.store.meta.deleteKeyIn(tx, key, n.namespace...)
		if err != nil {
			return fmt.Errorf("%v: %w", err, kv.ErrInternalError)
		}
		err = tx.Delete(k, &opt.WriteOptions{Sync: true})
		if err != nil {
			return fmt.Errorf("%v: %w", err, kv.ErrInternalError)
		}
		return nil
	})
}

// Iterate iterates the values in the namespace.
func (n *namespace) Iterate(
	ctx context.Context, f func(key string, value []byte) (next bool)) error {
	return n.store.read(func(db *leveldb.DB) error {
		keys, err := n.store.meta.getKeysIn(db, n.namespace...)
		if err != nil {
			return err
		}
		for _, key := range keys {
			k := n.store.meta.getKeyIn(key, n.namespace...)
			val, err := n.get(db, k)
			if err != nil && errors.Is(err, leveldb.ErrNotFound) {
				continue
			} else if err != nil {
				return err
			}
			if !f(key, val) {
				break
			}
		}
		return nil
	})
}

// Drop drops all the values in the namespace.
func (n *namespace) Drop(ctx context.Context) error {
	return n.store.write(n.drop)
}

func (n *namespace) drop(tx *leveldb.Transaction) error {
	// drop sub namespaces
	namespaces, err := n.store.meta.getNamespacesIn(tx, n.namespace...)
	if errors.Is(err, kv.ErrNotFound) {
		return nil
	} else if err != nil {
		return err
	}
	for _, nsName := range namespaces {
		subns := append(n.namespace[:0:0], n.namespace...)
		subns = append(subns, nsName)
		ns := &namespace{
			store:     n.store,
			namespace: subns,
		}
		err = ns.drop(tx)
		if err != nil {
			return err
		}
	}

	// drop values
	keys, err := n.store.meta.getKeysInTx(tx, n.namespace...)
	if err != nil {
		return err
	}
	for _, key := range keys {
		k := n.store.meta.getKeyIn(key, n.namespace...)
		err = tx.Delete(k, &opt.WriteOptions{Sync: true})
		if err != nil {
			return fmt.Errorf("%v: %w", err, kv.ErrInternalError)
		}
	}

	// drop meta data
	return n.store.meta.dropNamespaceMeta(tx, n.namespace...)
}

func (n *namespace) String() string {
	return fmt.Sprintf("leveldb.namespace(%s:%v)", n.store, n.namespace)
}
