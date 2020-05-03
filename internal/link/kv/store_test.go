package kv

import (
	"testing"

	"github.com/haostudio/golinks/internal/encoding/gob"
	"github.com/haostudio/golinks/internal/kv/memory"
	"github.com/haostudio/golinks/internal/link/linktest"
)

func TestStoreLogic(t *testing.T) {
	kvStore := memory.New()
	store := New(kvStore.In("test"), gob.New())
	linktest.StoreLogicTest(t, store)
}
