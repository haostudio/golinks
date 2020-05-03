package limiter

import "context"

// Limiter defines the limiter interface.
type Limiter interface {
	// TryHit tries to hit the limiter.
	TryHit(ctx context.Context, key string) (Result, error)
}

// Result defines the limiter result struct.
type Result struct {
	Limit     int
	Remaining int
	Reset     int64
	Reached   bool
}
