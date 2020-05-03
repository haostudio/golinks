package traced

import (
	"fmt"

	"github.com/haostudio/golinks/internal/kv"
)

// New returns a traced kv.Store.
func New(s kv.Store) kv.Store {
	return &store{
		store: s,
	}
}

type store struct {
	store kv.Store
}

// nolint: godox
// TODO: trace close function
// Close finalizes a kv.Store.
func (s *store) Close() error {
	return s.store.Close()
}

// In returns the namespace instance with path.
func (s *store) In(path ...string) kv.Namespace {
	return &namespace{
		store: s,
		root:  append(path[:0:0], path...),
		ns:    s.store.In(path...),
	}
}

func (s *store) String() string {
	return fmt.Sprintf("traced(%s)", s.store)
}
