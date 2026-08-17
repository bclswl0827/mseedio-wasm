package ms_doy2md

import (
	"context"
	"sync"

	"github.com/bclswl0827/mseedio-wasm/utils"
	"github.com/bclswl0827/mseedio-wasm/wrapper"
	"github.com/tetratelabs/wazero/api"
)

func New(ctx context.Context, mutex *sync.Mutex, function api.Function, valuePointerHandler utils.ValuePointerHandler) *IMPL_ms_doy2md {
	return &IMPL_ms_doy2md{
		BaseDependency:      wrapper.NewBaseDependency(ctx, mutex, function),
		valuePointerHandler: valuePointerHandler,
	}
}
