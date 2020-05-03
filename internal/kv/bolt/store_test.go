package bolt

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/haostudio/golinks/internal/kv/kvtest"
)

func TestLogic(t *testing.T) {
	// Prepare DB path
	dir, err := os.Getwd()
	require.NoError(t, err)
	dbPath := filepath.Join(dir,
		fmt.Sprintf("leveldb_test_%d", time.Now().UnixNano()))
	dbName := "test.db"
	// Test
	store, err := New(dbPath, dbName, "test")
	require.NoError(t, err)
	kvtest.StoreLogicTest(t, store)
	store, err = New(dbPath, dbName, "test")
	require.NoError(t, err)
	kvtest.StoreLogicTest(t, store)
	// Clean up
	require.NoError(t, os.RemoveAll(dbPath))
}

func TestConcurrent(t *testing.T) {
	// Prepare DB path
	dir, err := os.Getwd()
	require.NoError(t, err)
	dbPath := filepath.Join(dir,
		fmt.Sprintf("leveldb_test_%d", time.Now().UnixNano()))
	dbName := "test.db"
	store, err := New(dbPath, dbName, "test")
	require.NoError(t, err)
	// Test
	kvtest.StoreConcurrentTest(t, store, 1<<5, true)
	// Clean up
	require.NoError(t, os.RemoveAll(dbPath))
}
