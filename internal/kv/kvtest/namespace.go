package kvtest

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/haostudio/golinks/internal/kv"
)

// NamespaceLogicTest test the kv.Namespace logics.
func NamespaceLogicTest(t *testing.T, store kv.Namespace) {
	ctx := context.Background()
	defer require.NoError(t, store.Drop(ctx))

	namespace := "NAMESPACE"
	key := "KEY"
	value := []byte("VALUE")

	var (
		cloned []byte
		val    []byte
		err    error
	)

	// nothing found in both root and NAMESPACE
	_, err = store.In().Get(ctx, key)
	require.True(t, errors.Is(err, kv.ErrNotFound))
	_, err = store.In(namespace).Get(ctx, key)
	require.True(t, errors.Is(err, kv.ErrNotFound))

	// set in root
	cloned = append(value[:0:0], value...) // nolint: gocritic
	require.NoError(t, store.In().Set(ctx, key, cloned))
	val, err = store.In().Get(ctx, key)
	require.NoError(t, err)
	require.Equal(t, value, val)
	_, err = store.In(namespace).Get(ctx, key)
	require.True(t, errors.Is(err, kv.ErrNotFound))
	_, err = store.In(namespace).In(namespace).Get(ctx, key)
	require.True(t, errors.Is(err, kv.ErrNotFound))

	// set in NAMESPACE
	cloned = append(value[:0:0], value...) // nolint: gocritic
	require.NoError(t, store.In(namespace).Set(ctx, key, cloned))
	val, err = store.In(namespace).Get(ctx, key)
	require.NoError(t, err)
	require.Equal(t, value, val)
	_, err = store.In(namespace).In(namespace).Get(ctx, key)
	require.True(t, errors.Is(err, kv.ErrNotFound))

	// set in one level down
	cloned = append(value[:0:0], value...) // nolint: gocritic
	require.NoError(t,
		store.In(namespace).In(namespace).Set(ctx, key, cloned))
	val, err = store.In(namespace).In(namespace).Get(ctx, key)
	require.NoError(t, err)
	require.Equal(t, value, val)

	// delete in root
	require.NoError(t, store.In().Delete(ctx, key))
	_, err = store.In().Get(ctx, key)
	require.True(t, errors.Is(err, kv.ErrNotFound))
	val, err = store.In(namespace).Get(ctx, key)
	require.NoError(t, err)
	require.Equal(t, value, val)
	val, err = store.In(namespace).In(namespace).Get(ctx, key)
	require.NoError(t, err)
	require.Equal(t, value, val)

	// delete in NAMESPACE
	require.NoError(t, store.In(namespace).Delete(ctx, key))
	_, err = store.In().Get(ctx, key)
	require.True(t, errors.Is(err, kv.ErrNotFound))
	_, err = store.In(namespace).Get(ctx, key)
	require.True(t, errors.Is(err, kv.ErrNotFound))
	val, err = store.In(namespace).In(namespace).Get(ctx, key)
	require.NoError(t, err)
	require.Equal(t, value, val)

	// delete in NAMESPACE.NAMESPACE
	require.NoError(t,
		store.In(namespace).In(namespace).Delete(ctx, key))
	_, err = store.In().Get(ctx, key)
	require.True(t, errors.Is(err, kv.ErrNotFound))
	_, err = store.In(namespace).Get(ctx, key)
	require.True(t, errors.Is(err, kv.ErrNotFound))
	_, err = store.In(namespace).In(namespace).Get(ctx, key)
	require.True(t, errors.Is(err, kv.ErrNotFound))

	// set in all to test drop
	cloned = append(value[:0:0], value...) // nolint: gocritic
	require.NoError(t, store.In().Set(ctx, key, cloned))
	cloned = append(value[:0:0], value...) // nolint: gocritic
	require.NoError(t, store.In(namespace).Set(ctx, key, cloned))
	cloned = append(value[:0:0], value...) // nolint: gocritic
	require.NoError(t,
		store.In(namespace).In(namespace).Set(ctx, key, cloned))

	// check values
	val, err = store.In().Get(ctx, key)
	require.NoError(t, err)
	require.Equal(t, value, val)
	val, err = store.In(namespace).Get(ctx, key)
	require.NoError(t, err)
	require.Equal(t, value, val)
	val, err = store.In(namespace).In(namespace).Get(ctx, key)
	require.NoError(t, err)
	require.Equal(t, value, val)

	// drop in NAMESPACE
	require.NoError(t, store.In(namespace).Drop(ctx))
	val, err = store.In().Get(ctx, key)
	require.NoError(t, err)
	require.Equal(t, value, val)
	_, err = store.In(namespace).Get(ctx, key)
	require.True(t, errors.Is(err, kv.ErrNotFound))
	_, err = store.In(namespace).In(namespace).Get(ctx, key)
	require.True(t, errors.Is(err, kv.ErrNotFound))

	// set back all for root test
	cloned = append(value[:0:0], value...) // nolint: gocritic
	require.NoError(t, store.In().Set(ctx, key, cloned))
	cloned = append(value[:0:0], value...) // nolint: gocritic
	require.NoError(t, store.In(namespace).Set(ctx, key, cloned))
	cloned = append(value[:0:0], value...) // nolint: gocritic
	require.NoError(t,
		store.In(namespace).In(namespace).Set(ctx, key, cloned))

	// drop in root should drop everything
	require.NoError(t, store.In().Drop(ctx))
	_, err = store.In().Get(ctx, key)
	require.True(t, errors.Is(err, kv.ErrNotFound))
	_, err = store.In(namespace).Get(ctx, key)
	require.True(t, errors.Is(err, kv.ErrNotFound))
	_, err = store.In(namespace).In(namespace).Get(ctx, key)
	require.True(t, errors.Is(err, kv.ErrNotFound))
}

// NamespaceConcurrentTest test set/get/delete/iterate concurrently.
func NamespaceConcurrentTest(t *testing.T,
	store kv.Namespace, count int, iter bool) {
	ctx := context.Background()
	defer require.NoError(t, store.Drop(ctx))

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
			namespaceConcurrentTest(ctx, t, store, ns, val, count, iter)
			wg.Done()
		}(namespace, value)
	}
	wg.Wait()
}

func namespaceConcurrentTest(
	ctx context.Context, t *testing.T, store kv.Namespace,
	namespace []string, val string, count int, iter bool) {
	var wg sync.WaitGroup
	defer wg.Wait()

	wg.Add(3)

	// consider read more than write
	readCount := count << 1
	writeCount := count

	keyFunc := func(i int) string { return strconv.Itoa(i % 7) }
	validKey := func(k string) bool {
		i, err := strconv.Atoi(k)
		if err != nil {
			return false
		}
		return i >= 0 && i < 7
	}
	valFunc := func(k string) []byte { return []byte(k + val) }

	// get
	go func() {
		defer wg.Done()
		for i := 0; i < readCount; i++ {
			key := keyFunc(i)
			val, err := store.In(namespace...).Get(ctx, key)
			if err != nil {
				require.True(t, errors.Is(err, kv.ErrNotFound))
				continue
			}
			require.Equal(t, val, valFunc(key))
			// to make things worse
			val[0] = 2
		}
	}()

	// set
	go func() {
		defer wg.Done()
		for i := 0; i < writeCount; i++ {
			key := keyFunc(i)
			require.NoError(t, store.In(namespace...).Set(ctx, key, valFunc(key)))
		}
	}()

	// delete
	go func() {
		defer wg.Done()
		for i := 0; i < writeCount; i++ {
			require.NoError(t, store.In(namespace...).Delete(ctx, keyFunc(i)))
		}
	}()

	// iterate
	if !iter {
		return
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < readCount; i++ {
			err := store.In(namespace...).Iterate(ctx, func(k string, val []byte) bool {
				require.True(t, validKey(k), "invalid k", k)
				require.Equal(t, valFunc(k), val)
				// to make things worse
				val[0] = 1
				return true
			})
			if err != nil {
				require.True(t, errors.Is(err, kv.ErrNotFound))
			}
		}
	}()
}
