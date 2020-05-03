package bolt

import (
	"sync"
	"time"

	bolt "go.etcd.io/bbolt"
)

var singletons struct {
	sync.Mutex
	m map[string]*boltdb
}

type boltdb struct {
	db    *bolt.DB
	count int
}

// New returns a bolt kev-value store.
func getDB(file string) (db *bolt.DB, err error) {
	singletons.Lock()
	defer singletons.Unlock()

	if singletons.m == nil {
		singletons.m = make(map[string]*boltdb)
	}

	singleton, ok := singletons.m[file]
	if ok {
		singleton.count++
		db = singleton.db
		return
	}
	db, err = bolt.Open(
		file, 0666, &bolt.Options{Timeout: 10 * time.Second},
	)
	if err != nil {
		return nil, err
	}
	singletons.m[file] = &boltdb{db: db, count: 1}
	return
}

func close(file string) error {
	singletons.Lock()
	defer singletons.Unlock()
	singleton, ok := singletons.m[file]
	if !ok {
		return nil
	}
	var err error
	singleton.count--
	if singleton.count == 0 {
		err = singleton.db.Close()
		delete(singletons.m, file)
	}
	return err
}
