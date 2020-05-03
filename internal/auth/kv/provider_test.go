package kv

import (
	"testing"

	"github.com/haostudio/golinks/internal/auth/authtest"
	"github.com/haostudio/golinks/internal/encoding/gob"
	"github.com/haostudio/golinks/internal/kv/memory"
)

func TestLogic(t *testing.T) {
	provider := New(memory.New().In("test"), gob.New())
	authtest.ProviderLogicTest(t, provider)
}
