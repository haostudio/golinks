package slidingwindow

import (
	"context"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/haostudio/golinks/internal/limiter"
)

func TestSlidingWindow(t *testing.T) {
	var wg sync.WaitGroup
	keyCount := 100
	keys := make([]string, keyCount)
	for i := 0; i < keyCount; i++ {
		keys[i] = strconv.Itoa(i)
	}
	window := New(5, 6*time.Second)
	for _, key := range keys {
		wg.Add(1)
		go func(key string) {
			testSlidingWindow(t, window, key)
			wg.Done()
		}(key)
	}
	wg.Wait()
}

func testSlidingWindow(t *testing.T, window limiter.Limiter, key string) {
	ctx := context.Background()
	var res limiter.Result
	var err error
	for i := 0; i < 5; i++ {
		res, err = window.TryHit(ctx, key)
		require.NoError(t, err)
		require.False(t, res.Reached)
		time.Sleep(time.Second)
	}
	res, err = window.TryHit(ctx, key)
	require.NoError(t, err)
	require.True(t, res.Reached)
	time.Sleep(time.Second)

	res, err = window.TryHit(ctx, key)
	require.NoError(t, err)
	require.False(t, res.Reached)
	res, err = window.TryHit(ctx, key)
	require.NoError(t, err)
	require.True(t, res.Reached)
	time.Sleep(time.Second)

	res, err = window.TryHit(ctx, key)
	require.NoError(t, err)
	require.False(t, res.Reached)
	res, err = window.TryHit(ctx, key)
	require.NoError(t, err)
	require.True(t, res.Reached)
}
