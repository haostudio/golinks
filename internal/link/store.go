package link

import (
	"context"
	"fmt"
)

// Store defines the link store interface.
type Store interface {
	fmt.Stringer

	GetLink(ctx context.Context, org string, key string) (Link, error)
	GetLinks(ctx context.Context, org string) (map[string]Link, error)
	UpdateLink(ctx context.Context, org string, key string, ln Link) error
	DeleteLink(ctx context.Context, org string, key string) error
}
