package ms_mdtimestr2nstime

import (
	"fmt"

	types "github.com/bclswl0827/mseedio-wasm/wrapper/_types"
)

func (m *IMPL_ms_mdtimestr2nstime) Call(timeStr string) (int64, error) {
	if m.Ctx == nil {
		return 0, fmt.Errorf("context is nil")
	}
	if m.Mutex == nil {
		return 0, fmt.Errorf("mutex is nil")
	}
	if m.FunctionObj == nil {
		return 0, fmt.Errorf("ms_mdtimestr2nstime function is nil")
	}

	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	timeStrPtr, err := m.stringHandler.Pointer(timeStr)
	if err != nil {
		return 0, fmt.Errorf("failed to get pointer: %w", err)
	}
	defer m.stringHandler.Free(timeStrPtr)

	ret, err := m.FunctionObj.Call(m.Ctx, uint64(timeStrPtr))
	if err != nil {
		return 0, fmt.Errorf("failed to call ms_mdtimestr2nstime: %w", err)
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
