package memory

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/haostudio/golinks/internal/kv/kvtest"
)

func TestLRULogic(t *testing.T) {
	store := newStore(newLRUStore(10))
	defer func() {
		require.NoError(t, store.Close())
	}()
	kvtest.StoreLogicTest(t, store)
}

func TestLRUConcurrent(t *testing.T) {
	store := newStore(newLRUStore(10))
	kvtest.StoreConcurrentTest(t, store, 1<<12, false)
	require.NoError(t, store.Close())
}
