package traced

import (
	"context"
	"fmt"

	"go.opencensus.io/trace"

	"github.com/haostudio/golinks/internal/kv"
)

type namespace struct {
	store *store
	root  []string
	ns    kv.Namespace
}

// In returns the namespace instance with path.
func (n *namespace) In(path ...string) kv.Namespace {
	p := append(n.root[:0:0], n.root...)
	p = append(p, path...)
	return n.store.In(p...)
}

// Get returns the value in the namespace with key.
func (n *namespace) Get(ctx context.Context, key string) (b []byte, err error) {
	ctx, span := trace.StartSpan(ctx, "namespace.Get")
	defer span.End()
	span.AddAttributes(trace.StringAttribute("type", "kv_store"))
	span.AddAttributes(trace.StringAttribute("store", n.store.store.String()))
	span.AddAttributes(trace.StringAttribute("namespace", n.ns.String()))
	span.AddAttributes(trace.StringAttribute("kv_key", key))
	return n.ns.Get(ctx, key)
}

// Set sets the value in the namespace with key.
func (n *namespace) Set(ctx context.Context, key string, value []byte) (
	err error) {
	ctx, span := trace.StartSpan(ctx, "namespace.Set")
	defer span.End()
	span.AddAttributes(trace.StringAttribute("type", "kv_store"))
	span.AddAttributes(trace.StringAttribute("store", n.store.store.String()))
	span.AddAttributes(trace.StringAttribute("namespace", n.ns.String()))
	span.AddAttributes(trace.StringAttribute("kv_key", key))
	return n.ns.Set(ctx, key, value)
}

// Delete deletes the value in the namespace with key.
func (n *namespace) Delete(ctx context.Context, key string) (err error) {
	ctx, span := trace.StartSpan(ctx, "namespace.Delete")
	defer span.End()
	span.AddAttributes(trace.StringAttribute("type", "kv_store"))
	span.AddAttributes(trace.StringAttribute("store", n.store.store.String()))
	span.AddAttributes(trace.StringAttribute("namespace", n.ns.String()))
	span.AddAttributes(trace.StringAttribute("kv_key", key))
	return n.ns.Delete(ctx, key)
}

// Iterate iterates the values in the namespace.
func (n *namespace) Iterate(ctx context.Context,
	f func(key string, value []byte) (next bool)) (err error) {
	ctx, span := trace.StartSpan(ctx, "namespace.Iterate")
	defer span.End()
	span.AddAttributes(trace.StringAttribute("type", "kv_store"))
	span.AddAttributes(trace.StringAttribute("store", n.store.store.String()))
	span.AddAttributes(trace.StringAttribute("namespace", n.ns.String()))
	return n.ns.Iterate(ctx, f)
}

// Drop drops all the values in the namespace.
func (n *namespace) Drop(ctx context.Context) (err error) {
	ctx, span := trace.StartSpan(ctx, "namespace.Drop")
	defer span.End()
	span.AddAttributes(trace.StringAttribute("type", "kv_store"))
	span.AddAttributes(trace.StringAttribute("store", n.store.store.String()))
	span.AddAttributes(trace.StringAttribute("namespace", n.ns.String()))
	return n.ns.Drop(ctx)
}

func (n *namespace) String() string {
	return fmt.Sprintf("traced.namespace(%s)", n.ns)
}
