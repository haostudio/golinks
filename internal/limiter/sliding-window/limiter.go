package slidingwindow

import (
	"container/list"
	"context"
	"sync"
	"time"

	"github.com/haostudio/golinks/internal/limiter"
)

// New returns a new sliding window limiter.
func New(limit int, size time.Duration) limiter.Limiter {
	return &slidingWindow{
		limit:   limit,
		size:    size,
		windows: make(map[string]*window),
	}
}

type slidingWindow struct {
	sync.Mutex
	limit   int
	size    time.Duration
	windows map[string]*window
}

func (w *slidingWindow) TryHit(ctx context.Context, key string) (
	res limiter.Result, err error) {
	w.Lock()
	wd, ok := w.windows[key]
	if !ok {
		wd = &window{
			list: new(list.List),
		}
		w.windows[key] = wd
	}
	w.Unlock()
	res = wd.tryHit(w.limit, w.size)
	return
}

type window struct {
	sync.Mutex
	list *list.List
}

func (w *window) tryHit(limit int, size time.Duration) limiter.Result {
	w.Lock()
	defer w.Unlock()
	head := w.pruneAndPeek(size)

	// calculate reset time
	now := time.Now()
	reset := now.Add(size)
	if head != nil {
		reset = head.Value.(time.Time).Add(size)
	}

	// within limit
	l := w.list.Len()
	if l < limit {
		w.list.PushBack(now)
		return limiter.Result{
			Limit:     limit,
			Remaining: limit - l - 1,
			Reset:     reset.Unix(),
			Reached:   false,
		}
	}

	// over the limit
	return limiter.Result{
		Limit:     limit,
		Remaining: 0,
		Reset:     reset.Unix(),
		Reached:   true,
	}
}

func (w *window) pruneAndPeek(size time.Duration) *list.Element {
	l := w.list.Len()
	for l > 0 {
		elem := w.list.Front()
		expire := elem.Value.(time.Time)
		if expire.Add(size).Before(time.Now()) {
			// can abandon the elem
			w.list.Remove(elem)
			l--
			continue
		}
		return elem
	}
	return nil
}
