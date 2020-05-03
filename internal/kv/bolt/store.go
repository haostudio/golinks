package bolt

import (
	"fmt"
	"os"
	"path/filepath"

	bolt "go.etcd.io/bbolt"

	"github.com/haostudio/golinks/internal/kv"
)

type store struct {
	path, name, root string
	db               *bolt.DB
}

// New returns a bolt kev-value store.
func New(path, name, root string) (kv.Store, error) {
	// mkdir
	_, err := os.Stat(path)
	if err != nil {
		err = os.MkdirAll(path, os.ModePerm)
	}
	if err != nil {
		return nil, err
	}
	// open the database
	file := filepath.Join(path, name)
	db, err := getDB(file)
	if err != nil {
		return nil, err
	}
	s := &store{
		path: path,
		name: name,
		root: root,
		db:   db,
	}
	return s, nil
}

// Close finalizes a kv.Store.
func (s *store) Close() error {
	file := filepath.Join(s.path, s.name)
	return close(file)
}

// In returns the namespace instance with path.
func (s *store) In(path ...string) kv.Namespace {
	return &namespace{s, path}
}

func (s *store) String() string {
	return fmt.Sprintf("bolt.store(%s:%s:%s)", s.path, s.name, s.root)
}
