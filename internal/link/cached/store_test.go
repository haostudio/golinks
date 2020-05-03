package cached

import (
	"testing"

	"github.com/haostudio/golinks/internal/encoding/gob"
	"github.com/haostudio/golinks/internal/kv/memory"
	"github.com/haostudio/golinks/internal/link/kv"
	"github.com/haostudio/golinks/internal/link/linktest"
)

func TestStoreLogic(t *testing.T) {
	kvStore := memory.New()
	enc := gob.New()
	canonical := kv.New(kvStore.In("test"), enc)
	store := New(canonical, kvStore.In("cache"), enc)
	linktest.StoreLogicTest(t, store)
}
