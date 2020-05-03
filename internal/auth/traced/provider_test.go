package traced

import (
	"testing"

	"github.com/haostudio/golinks/internal/auth/authtest"
	"github.com/haostudio/golinks/internal/auth/kv"
	"github.com/haostudio/golinks/internal/encoding/gob"
	"github.com/haostudio/golinks/internal/kv/memory"
)

func TestLogic(t *testing.T) {
	provider := kv.New(memory.New().In("test"), gob.New())
	provider = New(provider)
	authtest.ProviderLogicTest(t, provider)
}
