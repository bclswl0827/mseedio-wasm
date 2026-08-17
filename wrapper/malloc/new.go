package malloc

import (
	"context"
	"sync"

	"github.com/bclswl0827/mseedio-wasm/wrapper"
	"github.com/tetratelabs/wazero/api"
)

func New(ctx context.Context, mutex *sync.Mutex, function api.Function) *IMPL_malloc {
	return &IMPL_malloc{
		BaseDependency: wrapper.NewBaseDependency(ctx, mutex, function),
	}
}
