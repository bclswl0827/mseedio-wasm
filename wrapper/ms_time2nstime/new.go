package ms_time2nstime

import (
	"context"
	"sync"

	"github.com/bclswl0827/mseedio-wasm/wrapper"
	"github.com/tetratelabs/wazero/api"
)

func New(ctx context.Context, mutex *sync.Mutex, function api.Function) *IMPL_ms_time2nstime {
	return &IMPL_ms_time2nstime{
		BaseDependency: wrapper.NewBaseDependency(ctx, mutex, function),
	}
}
