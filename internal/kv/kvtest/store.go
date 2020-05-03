package kvtest

import (
	"testing"

	"github.com/haostudio/golinks/internal/kv"
)

// StoreLogicTest test the kv.Store logics.
func StoreLogicTest(t *testing.T, store kv.Store) {
	NamespaceLogicTest(t, store.In())
}

// StoreConcurrentTest test set/get/delete/iterate concurrently.
func StoreConcurrentTest(t *testing.T, store kv.Store, count int, iter bool) {
	NamespaceConcurrentTest(t, store.In(), count, iter)
}
