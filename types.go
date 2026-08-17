package mseedio

import (
	"context"
	"sync"

	"github.com/bclswl0827/mseedio-wasm/wrapper/free"
	"github.com/bclswl0827/mseedio-wasm/wrapper/malloc"
	"github.com/bclswl0827/mseedio-wasm/wrapper/ms_doy2md"
	"github.com/bclswl0827/mseedio-wasm/wrapper/ms_md2doy"
	"github.com/bclswl0827/mseedio-wasm/wrapper/ms_mdtimestr2nstime"
	"github.com/bclswl0827/mseedio-wasm/wrapper/ms_nstime2time"
	"github.com/bclswl0827/mseedio-wasm/wrapper/ms_seedtimestr2nstime"
	"github.com/bclswl0827/mseedio-wasm/wrapper/ms_time2nstime"
	"github.com/bclswl0827/mseedio-wasm/wrapper/ms_timestr2nstime"
	"github.com/bclswl0827/mseedio-wasm/wrapper/msr3_init"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

type MiniSEED struct {
	ctx context.Context

	runtime wazero.Runtime
	module  api.Module
	memory  api.Memory

	mu          *sync.Mutex
	initialized bool

	// Raw exposes every function retained by the WASM export policy.
	Raw *RawAPI

	// Functions is retained for compatibility with the initial typed wrappers.
	// New code should prefer the high-level MiniSEED methods or Raw.
	Functions exportedFunction
}

type exportedFunction struct {
	FUNC_malloc *malloc.IMPL_malloc
	FUNC_free   *free.IMPL_free

	FUNC_ms_nstime2time        *ms_nstime2time.IMPL_ms_nstime2time
	FUNC_ms_time2nstime        *ms_time2nstime.IMPL_ms_time2nstime
	FUNC_ms_timestr2nstime     *ms_timestr2nstime.IMPL_ms_timestr2nstime
	FUNC_ms_mdtimestr2nstime   *ms_mdtimestr2nstime.IMPL_ms_mdtimestr2nstime
	FUNC_ms_seedtimestr2nstime *ms_seedtimestr2nstime.IMPL_ms_seedtimestr2nstime
	FUNC_ms_doy2md             *ms_doy2md.IMPL_ms_doy2md
	FUNC_ms_md2doy             *ms_md2doy.IMPL_ms_md2doy

	FUNC_msr3_init *msr3_init.IMPL_msr3_init
}
