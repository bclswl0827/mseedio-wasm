package free

import (
	"fmt"
)

func (f *IMPL_free) Call(ptr uint32) error {
	if f.Ctx == nil {
		return fmt.Errorf("context is nil")
	}
	if f.FunctionObj == nil {
		return fmt.Errorf("free function is nil")
	}
	if f.Mutex != nil {
		f.Mutex.Lock()
		defer f.Mutex.Unlock()
	}

	_, err := f.FunctionObj.Call(f.Ctx, uint64(ptr))
	if err != nil {
		return fmt.Errorf("failed to call free: %w", err)
	}

	return nil
}
