package memory

import (
	"container/list"
	"fmt"
	"strings"
	"sync"

	"github.com/haostudio/golinks/internal/kv"
)

const (
	lruStoreKeyPrefix    = "github.com/haostudio/golinks/internal/kv/memory/lru_"
	lruStoreKeySeparator = "/"
)

type lruStore struct {
	sync.RWMutex
	cap int                      // capacity
	l   *list.List               // doubly linked list
	m   map[string]*list.Element // hash table for checking if list node exists
}

type pair struct {
	key   string
	value []byte
}

func newLRUStore(capacity int) store {
	return &lruStore{
		cap: capacity,
		l:   new(list.List),
		m:   make(map[string]*list.Element, capacity),
	}
}

func (s *lruStore) key(keys ...string) string {
	elems := []string{lruStoreKeyPrefix}
	elems = append(elems, keys...)
	return strings.Join(elems, lruStoreKeySeparator)
}

func (s *lruStore) get(keys ...string) ([]byte, bool) {
	key := s.key(keys...)
	s.RLock()
	defer s.RUnlock()
	// check if list node exists
	if node, ok := s.m[key]; ok {
		val := node.Value.(pair).value
		// move node to front
		s.l.MoveToFront(node)
		cloned := append(val[:0:0], val...)
		return cloned, true
	}
	return nil, false
}

func (s *lruStore) set(value []byte, keys ...string) {
	key := s.key(keys...)
	s.Lock()
	defer s.Unlock()

	// update the existing one.
	if node, ok := s.m[key]; ok {
		s.l.MoveToFront(node)
		node.Value = pair{key: key, value: value}
		return
	}

	// delete the last list node if the list is full
	if s.l.Len() == s.cap {
		idx := s.l.Back().Value.(pair).key
		delete(s.m, idx)
		s.l.Remove(s.l.Back())
	}
	// push the new list node into the list
	ptr := s.l.PushFront(pair{
		key:   key,
		value: value,
	})
	s.m[key] = ptr
}

func (s *lruStore) del(keys ...string) {
	key := s.key(keys...)
	s.Lock()
	defer s.Unlock()
	node, ok := s.m[key]
	if !ok {
		return
	}
	s.l.Remove(node)
	delete(s.m, key)
}

func (s *lruStore) iter(
	f func(key string, val []byte) (next bool), key ...string,
) (exists bool, err error) {
	return false, kv.ErrNotSupport
}

func (s *lruStore) drop(keys ...string) {
	prefix := s.key(keys...) + lruStoreKeySeparator
	s.Lock()
	defer s.Unlock()
	for k, node := range s.m {
		if !strings.HasPrefix(k, prefix) {
			continue
		}
		s.l.Remove(node)
		delete(s.m, k)
	}
}
func (s *lruStore) String() string {
	return fmt.Sprintf("memory.lru-store(%p|%d)", s, s.cap)
}
