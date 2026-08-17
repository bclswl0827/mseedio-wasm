package ms_md2doy

import "fmt"

func (m *IMPL_ms_md2doy) Call(year, month, mday int32, yday *int32) (int, error) {
	if m.Ctx == nil {
		return 0, fmt.Errorf("context is nil")
	}
	if m.Mutex == nil {
		return 0, fmt.Errorf("mutex is nil")
	}
	if m.FunctionObj == nil {
		return 0, fmt.Errorf("ms_nstime2time function is nil")
	}

	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	ydayPtr := uint32(0)
	if yday != nil {
		ptr, err := m.valuePointerHandler.Pointer(*yday)
		if err != nil {
			return 0, fmt.Errorf("failed to get pointer of yday: %w", err)
		}
		ydayPtr = ptr
		defer m.valuePointerHandler.Free(ptr)
	}

	ret, err := m.FunctionObj.Call(m.Ctx, uint64(year), uint64(month), uint64(mday), uint64(ydayPtr))
	if err != nil {
		return 0, fmt.Errorf("failed to call ms_doy2md: %w", err)
	}

	if yday != nil {
		ydayVal := any(int32(0))
		if err := m.valuePointerHandler.Read(ydayPtr, &ydayVal); err != nil {
			return 0, fmt.Errorf("failed to read yday: %w", err)
		}
		*yday = ydayVal.(int32)
	}

	if ret[0] != 0 {
		code := int32(ret[0])
		return int(code), fmt.Errorf("library returns error code %d", code)
	}

	return 0, nil
}
