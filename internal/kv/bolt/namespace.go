package bolt

import (
	"context"
	"errors"
	"fmt"

	bolt "go.etcd.io/bbolt"

	"github.com/haostudio/golinks/internal/kv"
)

type namespace struct {
	store  *store
	bucket []string
}

// In returns the namespace instance with path.
func (n *namespace) In(path ...string) kv.Namespace {
	if len(path) == 0 {
		return n
	}
	b := append(n.bucket[:0:0], n.bucket...)
	b = append(b, path...)
	return &namespace{n.store, b}
}

// Get returns the value in the namespace with key.
func (n *namespace) Get(ctx context.Context, keyStr string) (
	b []byte, err error) {
	err = n.store.db.View(func(tx *bolt.Tx) error {
		root := n.root(tx)
		if root == nil {
			return kv.ErrNotFound
		}
		bucket := n.getBucket(root, n.bucket...)
		if bucket == nil {
			return kv.ErrNotFound
		}
		b = bucket.Get(key(keyStr))
		return nil
	})
	if err != nil {
		b = nil
		return
	}
	if b == nil {
		err = kv.ErrNotFound
	}
	return
}

// Set sets the value in the namespace with key.
func (n *namespace) Set(
	ctx context.Context, keyStr string, value []byte) error {
	return n.store.db.Update(func(tx *bolt.Tx) error {
		root, err := n.rootIfNotExists(tx)
		if err != nil {
			return err
		}
		bucket, err := n.getBucketIfNotExists(root, n.bucket...)
		if err != nil {
			return err
		}
		return bucket.Put(key(keyStr), value)
	})
}

// Delete deletes the value in the namespace with key.
func (n *namespace) Delete(ctx context.Context, keyStr string) error {
	return n.store.db.Update(func(tx *bolt.Tx) error {
		root := n.root(tx)
		if root == nil {
			return nil
		}
		bucket := n.getBucket(root, n.bucket...)
		if bucket == nil {
			return nil
		}
		return bucket.Delete(key(keyStr))
	})
}

// Iterate iterates the values in the namespace.
func (n *namespace) Iterate(
	ctx context.Context, f func(key string, value []byte) (next bool)) error {
	return n.store.db.View(func(tx *bolt.Tx) error {
		root := n.root(tx)
		if root == nil {
			return kv.ErrNotFound
		}
		bucket := n.getBucket(root, n.bucket...)
		if bucket == nil {
			return kv.ErrNotFound
		}
		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			next := f(string(k), v)
			if !next {
				return nil
			}
		}
		return nil
	})
}

// Drop drops all the values in the namespace.
func (n *namespace) Drop(ctx context.Context) error {
	return n.store.db.Update(func(tx *bolt.Tx) error {
		if len(n.bucket) == 0 {
			// drop everything
			err := tx.DeleteBucket(key(n.store.root))
			if errors.Is(err, bolt.ErrBucketNotFound) {
				return nil
			}
			return err
		}
		root := n.root(tx)
		if root == nil {
			return nil
		}

		// get the parent bucket of the target one to drop
		last := len(n.bucket) - 1
		bucket := n.getBucket(root, n.bucket[:last]...)
		if bucket == nil {
			return nil
		}

		err := bucket.DeleteBucket(key(n.bucket[last]))
		if errors.Is(err, bolt.ErrBucketNotFound) {
			return nil
		}
		return err
	})
}

func (n *namespace) root(tx *bolt.Tx) *bolt.Bucket {
	return tx.Bucket(key(n.store.root))
}

func (n *namespace) rootIfNotExists(tx *bolt.Tx) (*bolt.Bucket, error) {
	return tx.CreateBucketIfNotExists(key(n.store.root))
}

func (n *namespace) getBucket(
	bucket *bolt.Bucket, paths ...string) *bolt.Bucket {
	if len(paths) == 0 {
		return bucket
	}
	subBucket := bucket.Bucket(key(paths[0]))
	if subBucket == nil {
		return nil
	}
	return n.getBucket(subBucket, paths[1:]...)
}

func (n *namespace) getBucketIfNotExists(
	bucket *bolt.Bucket, paths ...string) (*bolt.Bucket, error) {
	if len(paths) == 0 {
		return bucket, nil
	}
	subBucket, err := bucket.CreateBucketIfNotExists(key(paths[0]))
	if err != nil {
		return nil, err
	}
	return n.getBucketIfNotExists(subBucket, paths[1:]...)
}

func (n *namespace) String() string {
	return fmt.Sprintf("bolt.namespace(%s:%v)", n.store, n.bucket)
}
