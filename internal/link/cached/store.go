package cached

import (
	"context"
	"fmt"
	"strings"

	"github.com/haostudio/golinks/internal/encoding"
	"github.com/haostudio/golinks/internal/kv"
	"github.com/haostudio/golinks/internal/link"
)

const (
	cachePrefix      = "github.com/haostudio/golinks/golinks/cached"
	allLinksCacheKey = "github.com/haostudio/golinks/golinks/cached_all"
)

// New returns a link.Store cached with a kv.Store.
func New(
	canonical link.Store, cache kv.Namespace, cacheEnc encoding.Binary,
) link.Store {
	s := &store{
		canonical: canonical,
	}
	s.cache.kv = cache
	s.cache.enc = cacheEnc
	return s
}

type store struct {
	canonical link.Store
	cache     struct {
		kv  kv.Namespace
		enc encoding.Binary
	}
}

func (s *store) GetLink(ctx context.Context, org string, key string) (
	ln link.Link, err error) {
	cacheKey := s.cacheKey(org, key)
	b, err := s.cache.kv.Get(ctx, cacheKey)
	for err == nil {
		err = s.cache.enc.Decode(b, &ln)
		if err != nil {
			// deal with decode error as cache miss and do the best effort to delete
			// dirty data.
			err = nil
			_ = s.cache.kv.Delete(ctx, cacheKey)
			break
		}
		return
	}
	ln, err = s.canonical.GetLink(ctx, org, key)
	if err != nil {
		return
	}
	for {
		// best effort to set the cache
		b, err = s.cache.enc.Encode(ln)
		if err != nil {
			// ignore this error
			err = nil
			break
		}
		_ = s.cache.kv.Set(ctx, cacheKey, b)
		break
	}
	return
}

func (s *store) GetLinks(ctx context.Context, org string) (
	links map[string]link.Link, err error) {
	b, err := s.cache.kv.Get(ctx, allLinksCacheKey)
	for err == nil {
		err = s.cache.enc.Decode(b, &links)
		if err != nil {
			// deal with decode error as cache miss and do the best effort to delete
			// dirty data.
			err = nil
			_ = s.cache.kv.Delete(ctx, allLinksCacheKey)
			break
		}
		return
	}
	links, err = s.canonical.GetLinks(ctx, org)
	if err != nil {
		return
	}
	for {
		// best effort to set the cache
		b, err = s.cache.enc.Encode(links)
		if err != nil {
			// ignore this error
			err = nil
			break
		}
		_ = s.cache.kv.Set(ctx, allLinksCacheKey, b)
		break
	}
	return
}

func (s *store) UpdateLink(
	ctx context.Context, org string, key string, ln link.Link) error {
	cacheKey := s.cacheKey(org, key)
	err := s.canonical.UpdateLink(ctx, org, key, ln)
	if err != nil {
		return err
	}
	// best effort to update cache
	b, err := s.cache.enc.Encode(ln)
	if err != nil {
		// ignore this error
		return nil
	}
	// ignore this error
	_ = s.cache.kv.Set(ctx, cacheKey, b)
	return nil
}

func (s *store) DeleteLink(ctx context.Context, org string, key string) error {
	cacheKey := s.cacheKey(org, key)
	err := s.canonical.DeleteLink(ctx, org, key)
	if err != nil {
		return err
	}
	// best effort to delete cache
	// ignore this error
	_ = s.cache.kv.Delete(ctx, cacheKey)
	return nil
}

func (s *store) cacheKey(org, key string) string {
	return strings.Join([]string{cachePrefix, org, key}, ".")
}

func (s *store) String() string {
	return fmt.Sprintf(
		"cacehd.store(%s/%s|%s)", s.cache.kv, s.cache.enc, s.canonical)
}
