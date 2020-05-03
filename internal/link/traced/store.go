package traced

import (
	"context"
	"fmt"

	"go.opencensus.io/trace"

	"github.com/haostudio/golinks/internal/link"
)

// New returns a traced link.Store.
func New(s link.Store) link.Store {
	return &store{
		store: s,
	}
}

type store struct {
	store link.Store
}

func (s *store) GetLink(ctx context.Context, org string, key string) (
	link.Link, error) {
	ctx, span := trace.StartSpan(ctx, "store.GetLink")
	defer span.End()
	span.AddAttributes(trace.StringAttribute("store", s.store.String()))
	return s.store.GetLink(ctx, org, key)
}

func (s *store) GetLinks(ctx context.Context, org string) (
	map[string]link.Link, error) {
	ctx, span := trace.StartSpan(ctx, "store.GetLinks")
	defer span.End()
	span.AddAttributes(trace.StringAttribute("store", s.store.String()))
	return s.store.GetLinks(ctx, org)
}

func (s *store) UpdateLink(
	ctx context.Context, org string, key string, ln link.Link) error {
	ctx, span := trace.StartSpan(ctx, "store.UpdateLink")
	defer span.End()
	span.AddAttributes(trace.StringAttribute("store", s.store.String()))
	return s.store.UpdateLink(ctx, org, key, ln)
}

func (s *store) DeleteLink(ctx context.Context, org string, key string) error {
	ctx, span := trace.StartSpan(ctx, "store.DeleteLink")
	defer span.End()
	span.AddAttributes(trace.StringAttribute("store", s.store.String()))
	return s.store.DeleteLink(ctx, org, key)
}

func (s *store) String() string {
	return fmt.Sprintf("traced(%s)", s.store)
}
