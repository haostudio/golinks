package memory

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/haostudio/golinks/internal/kv/kvtest"
)

func TestLogic(t *testing.T) {
	store := newStore(newMapStore())
	defer func() {
		require.NoError(t, store.Close())
	}()
	kvtest.StoreLogicTest(t, store)
	kvtest.StoreLogicTest(t, store)
}

func TestConcurrent(t *testing.T) {
	store := newStore(newMapStore())
	kvtest.StoreConcurrentTest(t, store, 1<<12, true)
	require.NoError(t, store.Close())
}
