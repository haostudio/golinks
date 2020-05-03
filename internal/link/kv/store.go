package kv

import (
	"context"
	"errors"
	"fmt"

	"github.com/haostudio/golinks/internal/encoding"
	"github.com/haostudio/golinks/internal/kv"
	"github.com/haostudio/golinks/internal/link"
)

// New returns a new store with kv and enc.
func New(kv kv.Namespace, enc encoding.Binary) link.Store {
	return &store{
		kv:  kv,
		enc: enc,
	}
}

type store struct {
	kv  kv.Namespace
	enc encoding.Binary
}

func (s *store) GetLink(ctx context.Context, org string, key string) (
	ln link.Link, err error) {
	// Get blob from kv
	b, err := s.kv.In(org).Get(ctx, key)
	if errors.Is(err, kv.ErrNotFound) {
		err = link.ErrNotFound
		return
	}
	if err != nil {
		return
	}
	// Decode
	err = s.enc.Decode(b, &ln)
	return
}

func (s *store) GetLinks(ctx context.Context, org string) (
	map[string]link.Link, error) {
	links := make(map[string]link.Link)
	err := s.kv.In(org).Iterate(ctx, func(key string, value []byte) bool {
		// Decode
		var ln link.Link
		iterErr := s.enc.Decode(value, &ln)
		if iterErr != nil {
			return true
		}
		links[key] = ln
		return true
	})
	if errors.Is(err, kv.ErrNotFound) {
		return nil, link.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return links, nil
}

func (s *store) UpdateLink(
	ctx context.Context, org string, key string, ln link.Link) error {
	blob, err := s.enc.Encode(ln)
	if err != nil {
		return err
	}
	return s.kv.In(org).Set(ctx, key, blob)
}

func (s *store) DeleteLink(ctx context.Context, org string, key string) error {
	return s.kv.In(org).Delete(ctx, key)
}

func (s *store) String() string {
	return fmt.Sprintf("kv.store(%s/%s)", s.kv, s.enc)
}
