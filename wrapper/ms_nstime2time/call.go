package ms_nstime2time

import (
	"fmt"
)

func (m *IMPL_ms_nstime2time) Call(nsTime int64, year, yday *uint16, hour, min, sec *uint8, nsec *uint32) (int64, error) {
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

	yearPtr := uint32(0)
	if year != nil {
		ptr, err := m.valuePointerHandler.Pointer(*year)
		if err != nil {
			return 0, fmt.Errorf("failed to get pointer of year: %w", err)
		}
		yearPtr = ptr
		defer m.valuePointerHandler.Free(ptr)
	}

	ydayPtr := uint32(0)
	if yday != nil {
		ptr, err := m.valuePointerHandler.Pointer(*yday)
		if err != nil {
			return 0, fmt.Errorf("failed to get pointer of yday: %w", err)
		}
		ydayPtr = ptr
		defer m.valuePointerHandler.Free(ptr)
	}

	hourPtr := uint32(0)
	if hour != nil {
		ptr, err := m.valuePointerHandler.Pointer(*hour)
		if err != nil {
			return 0, fmt.Errorf("failed to get pointer of hour: %w", err)
		}
		hourPtr = ptr
		defer m.valuePointerHandler.Free(ptr)
	}

	minPtr := uint32(0)
	if min != nil {
		ptr, err := m.valuePointerHandler.Pointer(*min)
		if err != nil {
			return 0, fmt.Errorf("failed to get pointer of min: %w", err)
		}
		minPtr = ptr
		defer m.valuePointerHandler.Free(ptr)
	}

	secPtr := uint32(0)
	if sec != nil {
		ptr, err := m.valuePointerHandler.Pointer(*sec)
		if err != nil {
			return 0, fmt.Errorf("failed to get pointer of sec: %w", err)
		}
		secPtr = ptr
		defer m.valuePointerHandler.Free(ptr)
	}

	nsecPtr := uint32(0)
	if nsec != nil {
		ptr, err := m.valuePointerHandler.Pointer(*nsec)
		if err != nil {
			return 0, fmt.Errorf("failed to get pointer of nsec: %w", err)
		}
		nsecPtr = ptr
		defer m.valuePointerHandler.Free(ptr)
	}

	ret, err := m.FunctionObj.Call(m.Ctx, uint64(nsTime), uint64(yearPtr), uint64(ydayPtr), uint64(hourPtr), uint64(minPtr), uint64(secPtr), uint64(nsecPtr))
	if err != nil {
		return 0, fmt.Errorf("failed to call ms_nstime2time: %w", err)
	}

	if year != nil {
		yearVal := any(uint16(0))
		if err := m.valuePointerHandler.Read(yearPtr, &yearVal); err != nil {
			return 0, fmt.Errorf("failed to read year: %w", err)
		}
		*year = yearVal.(uint16)
	}
	if yday != nil {
		ydayVal := any(uint16(0))
		if err := m.valuePointerHandler.Read(ydayPtr, &ydayVal); err != nil {
			return 0, fmt.Errorf("failed to read yday: %w", err)
		}
		*yday = ydayVal.(uint16)
	}
	if hour != nil {
		hourVal := any(uint8(0))
		if err := m.valuePointerHandler.Read(hourPtr, &hourVal); err != nil {
			return 0, fmt.Errorf("failed to read hour: %w", err)
		}
		*hour = hourVal.(uint8)
	}
	if min != nil {
		minVal := any(uint8(0))
		if err := m.valuePointerHandler.Read(minPtr, &minVal); err != nil {
			return 0, fmt.Errorf("failed to read min: %w", err)
		}
		*min = minVal.(uint8)
	}
	if sec != nil {
		secVal := any(uint8(0))
		if err := m.valuePointerHandler.Read(secPtr, &secVal); err != nil {
			return 0, fmt.Errorf("failed to read sec: %w", err)
		}
		*sec = secVal.(uint8)
	}
	if nsec != nil {
		nsecVal := any(uint32(0))
		if err := m.valuePointerHandler.Read(nsecPtr, &nsecVal); err != nil {
			return 0, fmt.Errorf("failed to read nsec: %w", err)
		}
		*nsec = nsecVal.(uint32)
	}

	return int64(int32(ret[0])), nil
}
