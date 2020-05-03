package linktest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/haostudio/golinks/internal/link"
)

// StoreLogicTest test the link.Store logics.
func StoreLogicTest(t *testing.T, store link.Store) {
	ctx := context.Background()

	org1 := "ORG_1"
	org2 := "ORG_2"
	key := "LINK"
	ln := link.V0("http://test")

	// nothing found in both root and NAMESPACE
	_, err := store.GetLink(ctx, org1, key)
	require.Error(t, link.ErrNotFound, err)
	_, err = store.GetLink(ctx, org2, key)
	require.Error(t, link.ErrNotFound, err)

	// set in root
	require.NoError(t, store.UpdateLink(ctx, org1, key, ln))
	l, err := store.GetLink(ctx, org1, key)
	require.NoError(t, err)
	require.Equal(t, ln, l)
	_, err = store.GetLink(ctx, org2, key)
	require.Error(t, link.ErrNotFound, err)

	// set in NAMESPACE
	require.NoError(t, store.UpdateLink(ctx, org2, key, ln))
	l, err = store.GetLink(ctx, org2, key)
	require.NoError(t, err)
	require.Equal(t, ln, l)

	// delete in root
	require.NoError(t, store.DeleteLink(ctx, org1, key))
	_, err = store.GetLink(ctx, org1, key)
	require.Error(t, link.ErrNotFound, err)
	l, err = store.GetLink(ctx, org2, key)
	require.NoError(t, err)
	require.Equal(t, ln, l)

	// delete in NAMESPACE
	require.NoError(t, store.DeleteLink(ctx, org2, key))
	_, err = store.GetLink(ctx, org1, key)
	require.Error(t, link.ErrNotFound, err)
	_, err = store.GetLink(ctx, org2, key)
	require.Error(t, link.ErrNotFound, err)
}

/*
// StoreConcurrentTest test set/get/delete/iterate concurrently.
func StoreConcurrentTest(t *testing.T, store link.Store, count int) {
	ctx := context.Background()

	var wg sync.WaitGroup

	// value, namespace
	cases := map[string][]string{
		"V":  []string{},
		"A":  []string{"A"},
		"AB": []string{"A", "B"},
		"BB": []string{"B", "B"},
	}

	for value, namespace := range cases {
		wg.Add(1)
		go func(ns []string, val string) {
			storeConcurrentTest(ctx, t, store, ns, "K", val, count)
			wg.Done()
		}(namespace, value)
	}
	wg.Wait()
}

func storeConcurrentTest(
	ctx context.Context, t *testing.T,
	store link.Store, namespace []string, key, val string, count int) {

	var wg sync.WaitGroup
	wg.Add(4)

	value := []byte(val)

	// consider read more than write
	readCount := count << 2
	writeCount := count

	// get
	go func() {
		defer wg.Done()
		for i := 0; i < readCount; i++ {
			val, err := store.In(namespace...).Get(ctx, key)
			if err != nil {
				require.True(t, errors.Is(err, link.ErrNotFound))
				continue
			}
			require.Equal(t, value, val)
			// to make things worse
			val[0] = 2
		}
	}()

	// set
	go func() {
		defer wg.Done()
		for i := 0; i < writeCount; i++ {
			cloned := make([]byte, len(value))
			copy(cloned, value)
			require.NoError(t, store.In(namespace...).Set(ctx, key, cloned))
		}
	}()

	// delete
	go func() {
		defer wg.Done()
		for i := 0; i < writeCount; i++ {
			require.NoError(t, store.In(namespace...).Delete(ctx, key))
		}
	}()

	// iterate
	go func() {
		defer wg.Done()
		for i := 0; i < readCount; i++ {
			err := store.In(namespace...).Iterate(ctx, func(k string, val []byte) bool {
				require.Equal(t, k, key)
				require.Equal(t, value, val)
				// to make things worse
				val[0] = 1
				return true
			})
			if err != nil {
				require.True(t, errors.Is(err, link.ErrNotFound))
			}
		}
	}()

	wg.Wait()
}
*/
