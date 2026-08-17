package msr3_init

import (
	"context"
	"sync"

	"github.com/bclswl0827/mseedio-wasm/utils"
	"github.com/bclswl0827/mseedio-wasm/wrapper"
	"github.com/tetratelabs/wazero/api"
)

func New(ctx context.Context, mutex *sync.Mutex, function api.Function, structHanlder utils.StructHandler) *IMPL_msr3_init {
	return &IMPL_msr3_init{
		BaseDependency: wrapper.NewBaseDependency(ctx, mutex, function),
		structHanlder:  structHanlder,
	}
}
