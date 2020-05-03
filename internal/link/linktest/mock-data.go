package linktest

import (
	"context"

	"github.com/haostudio/golinks/internal/encoding"
	"github.com/haostudio/golinks/internal/link"
)

// CreateSampleStore creates sample data in store with enc.
func CreateSampleStore(ctx context.Context,
	store link.Store, enc encoding.Binary, org string) {
	pr, err := link.V2("https://github.com/haostudio/{0}/issues/created_by/{1}")
	if err != nil {
		panic(err)
	}
	links := map[string]link.Link{
		"g":             link.V0("https://google.com"),
		"git":           link.V0("https://github.com"),
		"git.haostudio": link.V1("https://github.com/haostudio/{}"),
		"git.pr":        pr,
	}

	for k, l := range links {
		err := store.UpdateLink(ctx, org, k, l)
		if err != nil {
			panic(err)
		}
	}
}
