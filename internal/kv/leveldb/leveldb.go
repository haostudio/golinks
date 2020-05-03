package leveldb

import (
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
)

var singletons struct {
	sync.Mutex
	m map[string]*leveldbdb
}

type leveldbdb struct {
	db    *leveldb.DB
	count int
}

// New returns a leveldb kev-value store.
func getDB(file string) (db *leveldb.DB, err error) {
	singletons.Lock()
	defer singletons.Unlock()

	if singletons.m == nil {
		singletons.m = make(map[string]*leveldbdb)
	}

	singleton, ok := singletons.m[file]
	if ok {
		singleton.count++
		db = singleton.db
		return
	}
	db, err = leveldb.OpenFile(file, nil)
	if _, corrupted := err.(*errors.ErrCorrupted); corrupted {
		db, err = leveldb.RecoverFile(file, nil)
	}
	if err != nil {
		return
	}
	singletons.m[file] = &leveldbdb{db: db, count: 1}
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
