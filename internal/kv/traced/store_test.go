package traced

import (
	"testing"

	"github.com/haostudio/golinks/internal/kv/kvtest"
	"github.com/haostudio/golinks/internal/kv/memory"
)

func TestLogic(t *testing.T) {
	store := New(memory.New())
	kvtest.StoreLogicTest(t, store)
}

func TestConcurrent(t *testing.T) {
	store := New(memory.New())
	kvtest.StoreConcurrentTest(t, store, 1<<12, false)
}
