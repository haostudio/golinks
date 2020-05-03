package cached

import (
	"fmt"
	"strings"
	"sync"

	"github.com/haostudio/golinks/internal/kv"
)

const keyPrefix = "github.com/haostudio/golinks/internal/kv/cached_namespaces"

// New returns a cache kv.Store.
func New(canonical kv.Store, cache kv.Namespace) kv.Store {
	return &store{
		canonical: canonical,
		cache:     cache,
	}
}

type store struct {
	canonical kv.Store
	cache     kv.Namespace

	namespaces struct {
		sync.Mutex
		m map[string]*namespace
	}
}

// Close finalizes a kv.Store.
func (s *store) Close() error {
	return kv.ErrNotSupport
}

// In returns the namespace instance with path.
func (s *store) In(path ...string) kv.Namespace {
	elem := append([]string{keyPrefix}, path...)
	key := strings.Join(elem, "_")
	s.namespaces.Lock()
	defer s.namespaces.Unlock()

	if s.namespaces.m == nil {
		s.namespaces.m = make(map[string]*namespace)
	}

	singleton, ok := s.namespaces.m[key]
	if ok {
		return singleton
	}

	s.namespaces.m[key] = &namespace{
		store:     s,
		root:      append(path[:0:0], path...),
		canonical: s.canonical.In(append(path[:0:0], path...)...),
		cache:     s.cache.In(append(path[:0:0], path...)...),
	}
	return s.namespaces.m[key]
}

func (s *store) String() string {
	return fmt.Sprintf("cached.store(%s/%s)", s.canonical, s.cache)
}
