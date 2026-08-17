package malloc

import (
	"fmt"
)

func (m *IMPL_malloc) Call(size uint32) (uint32, error) {
	if m.Ctx == nil {
		return 0, fmt.Errorf("context is nil")
	}
	if m.FunctionObj == nil {
		return 0, fmt.Errorf("malloc function is nil")
	}
	if m.Mutex != nil {
		m.Mutex.Lock()
		defer m.Mutex.Unlock()
	}

	ptr, err := m.FunctionObj.Call(m.Ctx, uint64(size))
	if err != nil {
		return 0, fmt.Errorf("failed to call malloc: %w", err)
	}

	if ptr[0] == 0 {
		return 0, fmt.Errorf("failed to allocate memory")
	}

	return uint32(ptr[0]), nil
}
