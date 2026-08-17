package mseedio

import (
	"fmt"
	"io"
	"os"

	"github.com/bclswl0827/mseedio-wasm/lib"
	"github.com/bclswl0827/mseedio-wasm/utils"
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
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

// InitOptions controls WASI streams and optional guest filesystem mounts.
type InitOptions struct {
	Stdout   io.Writer
	Stderr   io.Writer
	FSConfig wazero.FSConfig
}

// writerOnly prevents wazero from treating an *os.File as a WASI host file.
// Some terminals do not support Stat and return EINVAL during instantiation;
// the libmseed module only needs fd_write for these streams.
type writerOnly struct {
	io.Writer
}

// Init instantiates libmseed with standard output streams and no host mounts.
func (m *MiniSEED) Init(stdout, stderr *os.File) error {
	return m.InitWithOptions(InitOptions{Stdout: stdout, Stderr: stderr})
}

// InitWithOptions instantiates libmseed. FSConfig is required only when raw
// file-based libmseed functions need access to a host directory.
func (m *MiniSEED) InitWithOptions(options InitOptions) error {
	if m == nil || m.runtime == nil {
		return fmt.Errorf("mseedio: runtime is nil; create it with New")
	}
	if m.initialized {
		return fmt.Errorf("mseedio: runtime is already initialized")
	}
	if _, err := wasi_snapshot_preview1.Instantiate(m.ctx, m.runtime); err != nil {
		return fmt.Errorf("mseedio: instantiate WASI: %w", err)
	}

	config := wazero.NewModuleConfig()
	if options.Stdout == nil {
		options.Stdout = os.Stdout
	}
	if options.Stderr == nil {
		options.Stderr = os.Stderr
	}
	config = config.
		WithStdout(writerOnly{Writer: options.Stdout}).
		WithStderr(writerOnly{Writer: options.Stderr})
	if options.FSConfig != nil {
		config = config.WithFSConfig(options.FSConfig)
	}

	module, err := m.runtime.InstantiateWithConfig(m.ctx, lib.LoadBytes(), config)
	if err != nil {
		return fmt.Errorf("mseedio: instantiate libmseed: %w", err)
	}
	m.module = module
	m.memory = module.Memory()

	mallocFunction := module.ExportedFunction("malloc")
	freeFunction := module.ExportedFunction("free")
	mallocFn := malloc.New(m.ctx, m.mu, mallocFunction)
	freeFn := free.New(m.ctx, m.mu, freeFunction)
	internalMallocFn := malloc.New(m.ctx, nil, mallocFunction)
	internalFreeFn := free.New(m.ctx, nil, freeFunction)

	stringHandler := utils.NewStringHandler(m.memory, internalMallocFn, internalFreeFn)
	valuePointerHandler := utils.NewValuePointerHandler(m.memory, internalMallocFn, internalFreeFn)
	structHandler := utils.NewStructHandler(m.memory, internalMallocFn, internalFreeFn)

	m.Functions = exportedFunction{
		FUNC_malloc:                mallocFn,
		FUNC_free:                  freeFn,
		FUNC_ms_nstime2time:        ms_nstime2time.New(m.ctx, m.mu, module.ExportedFunction("ms_nstime2time"), valuePointerHandler),
		FUNC_ms_time2nstime:        ms_time2nstime.New(m.ctx, m.mu, module.ExportedFunction("ms_time2nstime")),
		FUNC_ms_timestr2nstime:     ms_timestr2nstime.New(m.ctx, m.mu, module.ExportedFunction("ms_timestr2nstime"), stringHandler),
		FUNC_ms_mdtimestr2nstime:   ms_mdtimestr2nstime.New(m.ctx, m.mu, module.ExportedFunction("ms_mdtimestr2nstime"), stringHandler),
		FUNC_ms_seedtimestr2nstime: ms_seedtimestr2nstime.New(m.ctx, m.mu, module.ExportedFunction("ms_seedtimestr2nstime"), stringHandler),
		FUNC_ms_doy2md:             ms_doy2md.New(m.ctx, m.mu, module.ExportedFunction("ms_doy2md"), valuePointerHandler),
		FUNC_ms_md2doy:             ms_md2doy.New(m.ctx, m.mu, module.ExportedFunction("ms_md2doy"), valuePointerHandler),

		FUNC_msr3_init: msr3_init.New(m.ctx, m.mu, module.ExportedFunction("msr3_init"), structHandler),
	}

	raw, err := newRawAPI(m)
	if err != nil {
		_ = module.Close(m.ctx)
		m.module = nil
		m.memory = nil
		m.Functions = exportedFunction{}
		return err
	}
	m.Raw = raw
	m.initialized = true
	return nil
}
