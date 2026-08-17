package ms_time2nstime

import (
	"fmt"

	types "github.com/bclswl0827/mseedio-wasm/wrapper/_types"
)

func (m *IMPL_ms_time2nstime) Call(year, yday, hour, min, sec int, nsec uint32) (int64, error) {
	if m.Ctx == nil {
		return 0, fmt.Errorf("context is nil")
	}
	if m.Mutex == nil {
		return 0, fmt.Errorf("mutex is nil")
	}
	if m.FunctionObj == nil {
		return 0, fmt.Errorf("ms_timestr2nstime function is nil")
	}

	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	ret, err := m.FunctionObj.Call(m.Ctx, uint64(year), uint64(yday), uint64(hour), uint64(min), uint64(sec), uint64(nsec))
	if err != nil {
		return 0, fmt.Errorf("failed to call ms_time2nstime: %w", err)
	}

	result := int64(ret[0])
	switch result {
	case types.NSTERROR:
		return result, fmt.Errorf("library returns error code %d", result)
	case types.NSTUNSET:
		return result, fmt.Errorf("library returns error code %d", result)
	}

	return result, nil
}
