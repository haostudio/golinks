package leveldb

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"

	"github.com/haostudio/golinks/internal/kv"
)

type store struct {
	path string
	name string
	meta *meta
	// XXX: This should no be necessary since we are using db transaction.
	// However, the concurrent consistency tests don't pass unless we do
	// this code level read/write lock.
	// XXX: It doesn't pass in this way either :(. run with -short at the moment.
	io struct {
		sync.RWMutex
		db *leveldb.DB
	}
}

// New returns a level kev-value store.
func New(path, name string, metaCache kv.Store) (kv.Store, error) {
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
		meta: &meta{},
	}
	s.io.db = db
	return s, nil
}

// Close finalizes a kv.Store.
func (s *store) Close() error {
	s.io.Lock()
	defer s.io.Unlock()
	file := filepath.Join(s.path, s.name)
	return close(file)
}

// In returns the namespace instance with path.
func (s *store) In(path ...string) kv.Namespace {
	return &namespace{
		store:     s,
		namespace: path,
	}
}

func (s *store) read(f func(db *leveldb.DB) error) (err error) {
	s.io.RLock()
	defer s.io.RUnlock()
	// s.io.Lock()
	// defer s.io.Unlock()
	return f(s.io.db)
}

func (s *store) write(f func(tx *leveldb.Transaction) error) (err error) {
	s.io.Lock()
	defer s.io.Unlock()
	tx, err := s.io.db.OpenTransaction()
	if err != nil {
		return
	}

	// Error checking and panic safenet.
	defer func() {
		if recovered := recover(); recovered != nil || err != nil {
			tx.Discard()
			if recovered != nil {
				panic(recovered)
			}
		}
	}()

	err = f(tx)
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}

func (s *store) String() string {
	return fmt.Sprintf("leveldb.store(%s:%s)", s.path, s.name)
}
