package ms_timestr2nstime

import (
	"context"
	"sync"

	"github.com/bclswl0827/mseedio-wasm/utils"
	"github.com/bclswl0827/mseedio-wasm/wrapper"
	"github.com/tetratelabs/wazero/api"
)

func New(ctx context.Context, mutex *sync.Mutex, function api.Function, stringHandler utils.StringHandler) *IMPL_ms_timestr2nstime {
	return &IMPL_ms_timestr2nstime{
		BaseDependency: wrapper.NewBaseDependency(ctx, mutex, function),
		stringHandler:  stringHandler,
	}
}
