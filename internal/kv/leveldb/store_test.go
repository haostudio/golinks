package leveldb

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"

	"github.com/haostudio/golinks/internal/kv/kvtest"
)

const testdb = "test.db"

func TestLogic(t *testing.T) {
	// Prepare DB path
	dir, err := os.Getwd()
	require.NoError(t, err)
	dbPath := filepath.Join(
		dir, fmt.Sprintf("leveldb_test_%d", time.Now().UnixNano()))
	dbName := testdb
	// Test twice
	store, err := New(dbPath, dbName, nil)
	defer func() {
		require.NoError(t, store.Close())
	}()
	require.NoError(t, err)
	kvtest.StoreLogicTest(t, store)
	store, err = New(dbPath, dbName, nil)
	defer func() {
		require.NoError(t, store.Close())
	}()
	require.NoError(t, err)
	kvtest.StoreLogicTest(t, store)
	// Clean up
	require.NoError(t, os.RemoveAll(dbPath))
}

func TestTransaction(t *testing.T) {
	// Prepare DB path
	dir, err := os.Getwd()
	require.NoError(t, err)
	dbPath := filepath.Join(
		dir, fmt.Sprintf("leveldb_test_%d", time.Now().UnixNano()))
	dbName := testdb
	kvStore, err := New(dbPath, dbName, nil)
	defer func() {
		require.NoError(t, kvStore.Close())
	}()
	require.NoError(t, err)
	leveldbStore, ok := kvStore.(*store)
	require.True(t, ok)
	// test failed tx
	xerr := errors.New("X")
	err = leveldbStore.write(func(tx *leveldb.Transaction) error {
		txErr := tx.Put([]byte("K"), []byte("V"), &opt.WriteOptions{Sync: true})
		require.NoError(t, txErr)
		return xerr
	})
	require.True(t, errors.Is(err, xerr))
	_, err = leveldbStore.io.db.Get([]byte("K"), nil)
	require.True(t, errors.Is(err, leveldb.ErrNotFound))
	// test success tx
	err = leveldbStore.write(func(tx *leveldb.Transaction) error {
		txErr := tx.Put([]byte("K"), []byte("V"), &opt.WriteOptions{Sync: true})
		require.NoError(t, txErr)
		return nil
	})
	require.NoError(t, err)
	b, err := leveldbStore.io.db.Get([]byte("K"), nil)
	require.NoError(t, err)
	require.Equal(t, []byte("V"), b)

	// Clean up
	require.NoError(t, os.RemoveAll(dbPath))
}

func TestConcurrentConsistency_1(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}
	// Prepare DB path
	dir, err := os.Getwd()
	require.NoError(t, err)
	dbPath := filepath.Join(
		dir, fmt.Sprintf("leveldb_test_%d", time.Now().UnixNano()))
	dbName := testdb
	kvStore, err := New(dbPath, dbName, nil)
	defer func() {
		require.NoError(t, kvStore.Close())
	}()
	require.NoError(t, err)
	leveldbStore, ok := kvStore.(*store)
	require.True(t, ok)
	// Test
	kvtest.StoreConcurrentTest(t, kvStore, 1<<5, true)

	// Check race condition by checking consistency between meta and actual slot.
	// XXX cases should be consistent with kv/kvtest/store.go
	key := "K"
	// value -> namespace
	cases := map[string][]string{
		"V":  []string{},
		"A":  []string{"A"},
		"AB": []string{"A", "B"},
		"BB": []string{"B", "B"},
	}
	// snapshot, err := leveldbStore.io.db.GetSnapshot()
	// require.NoError(t, err)
	// snapshot := leveldbStore.io.db
	require.NoError(t, leveldbStore.read(func(snapshot *leveldb.DB) error {
		for value, namespace := range cases {
			nsMetaKey := leveldbStore.meta.getNamespaceMetaKey(namespace...)
			valueKey := leveldbStore.meta.getKeyIn(key, namespace...)
			var shouldKeyExists bool
			b, err := snapshot.Get(nsMetaKey, nil)
			for {
				if errors.Is(err, leveldb.ErrNotFound) {
					break
				}
				require.NoError(t, err)
				var nsMeta namespaceMeta
				require.NoError(t, metaEnc.Decode(b, &nsMeta))
				fmt.Println(nsMeta)
				if nsMeta.Keys == nil {
					break
				}
				_, shouldKeyExists = nsMeta.Keys[key]
				break
			}
			b, err = snapshot.Get(valueKey, nil)
			if shouldKeyExists {
				require.NoError(t, err)
				require.Equal(t, []byte(value), b)
			} else {
				require.Error(t, err, string(b))
				require.True(t, errors.Is(err, leveldb.ErrNotFound))
			}
		}
		return nil
	}))
	// snapshot.Release()

	// Clean up
	require.NoError(t, os.RemoveAll(dbPath))
}

func TestConcurrentConsistency_2(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}
	// Prepare DB path
	dir, err := os.Getwd()
	require.NoError(t, err)
	dbPath := filepath.Join(
		dir, fmt.Sprintf("leveldb_test_%d", time.Now().UnixNano()))
	dbName := testdb
	kvStore, err := New(dbPath, dbName, nil)
	defer func() {
		require.NoError(t, kvStore.Close())
	}()
	require.NoError(t, err)
	leveldbStore, ok := kvStore.(*store)
	require.True(t, ok)
	// Test
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		kvtest.StoreConcurrentTest(t, kvStore, 1<<5, true)
		cancel()
	}()
	// Check race condition by checking consistency between meta and actual slot.
	// XXX cases should be consistent with kv/kvtest/store.go
	key := "K"
	// value -> namespace
	cases := map[string][]string{
		"V":  []string{},
		"A":  []string{"A"},
		"AB": []string{"A", "B"},
		"BB": []string{"B", "B"},
	}
	for ctx.Err() == nil {
		// snapshot, err := leveldbStore.io.db.GetSnapshot()
		// require.NoError(t, err)
		require.NoError(t, leveldbStore.read(func(snapshot *leveldb.DB) error {
			for value, namespace := range cases {
				nsMetaKey := leveldbStore.meta.getNamespaceMetaKey(namespace...)
				valueKey := leveldbStore.meta.getKeyIn(key, namespace...)
				var shouldKeyExists bool
				b, err := snapshot.Get(nsMetaKey, nil)
				for { // nolint: staticcheck
					if errors.Is(err, leveldb.ErrNotFound) {
						break
					}
					require.NoError(t, err)
					var nsMeta namespaceMeta
					require.NoError(t, metaEnc.Decode(b, &nsMeta))
					if nsMeta.Keys == nil {
						break
					}
					_, shouldKeyExists = nsMeta.Keys[key]
					break
				}
				b, err = snapshot.Get(valueKey, nil)
				if shouldKeyExists {
					require.NoError(t, err)
					require.Equal(t, []byte(value), b)
				} else {
					require.Error(t, err, string(b))
					require.True(t, errors.Is(err, leveldb.ErrNotFound))
				}
			}
			// snapshot.Release()
			return nil
		}))
		require.NoError(t, err)
	}

	// Clean up
	require.NoError(t, os.RemoveAll(dbPath))
}
