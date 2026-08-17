package mseedio

import "fmt"

func (m *MiniSEED) Close() error {
	if m == nil || m.runtime == nil {
		return nil
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.module != nil {
		if err := m.module.Close(m.ctx); err != nil {
			return fmt.Errorf("mseedio: close module: %w", err)
		}
		m.module = nil
	}
	if err := m.runtime.Close(m.ctx); err != nil {
		return fmt.Errorf("mseedio: close runtime: %w", err)
	}

	m.runtime = nil
	m.memory = nil
	m.Raw = nil
	m.Functions = exportedFunction{}
	m.initialized = false

	return nil
}
