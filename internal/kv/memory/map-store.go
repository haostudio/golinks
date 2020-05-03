package memory

import (
	"fmt"
	"sync"
)

type entry struct {
	value     []byte
	namespace map[string]*entry
}

type mapStore struct {
	sync.RWMutex
	m map[string]*entry
}

func newMapStore() store {
	return &mapStore{
		m: make(map[string]*entry),
	}
}

func (s *mapStore) get(keys ...string) (b []byte, exists bool) {
	s.RLock()
	defer s.RUnlock()
	v, exists := s.getNoLock(s.m, keys...)
	if !exists {
		return
	}
	b = append(v[:0:0], v...) // nolint: gocritic
	return b, true
}

func (s *mapStore) getNoLock(m map[string]*entry, keys ...string) (
	[]byte, bool) {
	e, ok := m[keys[0]]
	if !ok {
		return nil, false
	}
	if len(keys) == 1 {
		if e.value == nil {
			return nil, false
		}
		return e.value, true
	}
	if e.namespace == nil {
		return nil, false
	}
	return s.getNoLock(e.namespace, keys[1:]...)
}

func (s *mapStore) set(value []byte, keys ...string) {
	s.Lock()
	defer s.Unlock()
	s.setNoLock(s.m, value, keys...)
}

func (s *mapStore) setNoLock(
	m map[string]*entry, value []byte, keys ...string) {
	_, ok := m[keys[0]]
	if !ok {
		m[keys[0]] = &entry{}
	}

	if len(keys) == 1 {
		m[keys[0]].value = value
		return
	}

	if m[keys[0]].namespace == nil {
		m[keys[0]].namespace = make(map[string]*entry)
	}
	s.setNoLock(m[keys[0]].namespace, value, keys[1:]...)
}

func (s *mapStore) del(keys ...string) {
	s.Lock()
	defer s.Unlock()
	s.delNoLock(s.m, keys...)
}

func (s *mapStore) delNoLock(m map[string]*entry, keys ...string) {
	_, ok := m[keys[0]]
	if !ok {
		return
	}
	if len(keys) == 1 {
		m[keys[0]].value = nil
		return
	}
	if m[keys[0]].namespace == nil {
		return
	}
	s.delNoLock(m[keys[0]].namespace, keys[1:]...)
}

func (s *mapStore) iter(
	f func(key string, val []byte) (next bool), path ...string) (bool, error) {
	s.RLock()
	m, exists := s.getMapNoLock(s.m, path...)
	s.RUnlock()
	if !exists {
		return false, nil
	}
	for key, value := range m {
		if !f(key, value) {
			break
		}
	}
	return true, nil
}

func (s *mapStore) getMapNoLock(m map[string]*entry, keys ...string) (
	map[string][]byte, bool) {
	if len(keys) == 0 {
		clonedMap := make(map[string][]byte)
		for k, e := range m {
			if e.value == nil {
				continue
			}
			cloned := make([]byte, len(e.value))
			copy(cloned, e.value)
			clonedMap[k] = cloned
		}
		return clonedMap, true
	}
	e, ok := m[keys[0]]
	if !ok {
		return nil, false
	}
	if e.namespace == nil {
		return nil, false
	}
	return s.getMapNoLock(e.namespace, keys[1:]...)
}

func (s *mapStore) drop(path ...string) {
	s.Lock()
	defer s.Unlock()
	if len(path) == 0 {
		// drop everything
		s.m = make(map[string]*entry)
		return
	}
	s.dropNoLock(s.m, path...)
}

func (s *mapStore) dropNoLock(m map[string]*entry, path ...string) {
	e, ok := m[path[0]]
	if !ok {
		return
	}
	if len(path) == 1 {
		// drop the namespace
		e.namespace = nil
		// drop the entry if no value stored
		if e.value == nil {
			delete(m, path[0])
		}
		return
	}
	if e.namespace == nil {
		return
	}
	s.dropNoLock(e.namespace, path[1:]...)
}

func (s *mapStore) String() string {
	return fmt.Sprintf("memory.map-store(%p)", s)
}
