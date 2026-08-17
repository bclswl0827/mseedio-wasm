package ms_doy2md

import "fmt"

func (m *IMPL_ms_doy2md) Call(year, yday int32, month, mday *int32) (int, error) {
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

	monthPtr := uint32(0)
	if month != nil {
		ptr, err := m.valuePointerHandler.Pointer(*month)
		if err != nil {
			return 0, fmt.Errorf("failed to get pointer of month: %w", err)
		}
		monthPtr = ptr
		defer m.valuePointerHandler.Free(ptr)
	}

	mdayPtr := uint32(0)
	if mday != nil {
		ptr, err := m.valuePointerHandler.Pointer(*mday)
		if err != nil {
			return 0, fmt.Errorf("failed to get pointer of mday: %w", err)
		}
		mdayPtr = ptr
		defer m.valuePointerHandler.Free(ptr)
	}

	ret, err := m.FunctionObj.Call(m.Ctx, uint64(year), uint64(yday), uint64(monthPtr), uint64(mdayPtr))
	if err != nil {
		return 0, fmt.Errorf("failed to call ms_doy2md: %w", err)
	}

	if month != nil {
		monthVal := any(int32(0))
		if err := m.valuePointerHandler.Read(monthPtr, &monthVal); err != nil {
			return 0, fmt.Errorf("failed to read month: %w", err)
		}
		*month = monthVal.(int32)
	}

	if mday != nil {
		mdayVal := any(int32(0))
		if err := m.valuePointerHandler.Read(mdayPtr, &mdayVal); err != nil {
			return 0, fmt.Errorf("failed to read mday: %w", err)
		}
		*mday = mdayVal.(int32)
	}

	if ret[0] != 0 {
		code := int32(ret[0])
		return int(code), fmt.Errorf("library returns error code %d", code)
	}

	return 0, nil
}
