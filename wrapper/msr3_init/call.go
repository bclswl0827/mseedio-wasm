package msr3_init

import (
	"fmt"

	types "github.com/bclswl0827/mseedio-wasm/wrapper/_types"
)

func (m *IMPL_msr3_init) Call() (*types.M3Record, error) {
	if m.Ctx == nil {
		return nil, fmt.Errorf("context is nil")
	}
	if m.Mutex == nil {
		return nil, fmt.Errorf("mutex is nil")
	}
	if m.FunctionObj == nil {
		return nil, fmt.Errorf("msr3_init function is nil")
	}

	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	m3RecordAddr, err := m.structHanlder.Pointer(types.SizeOfM3Record)
	if err != nil {
		return nil, err
	}
	defer m.structHanlder.Free(m3RecordAddr)

	ret, err := m.FunctionObj.Call(m.Ctx, uint64(m3RecordAddr))
	if err != nil {
		return nil, err
	}

	retPtr := uint32(ret[0])
	if retPtr == 0 {
		return nil, fmt.Errorf("msr3_init returns nil")
	}

	m3RecordBytes, err := m.structHanlder.Read(retPtr, types.SizeOfM3Record)
	if err != nil {
		return nil, err
	}

	var m3Record types.M3Record
	if err = m3Record.FromBytes(m.structHanlder, m3RecordBytes); err != nil {
		return nil, err
	}

	return &m3Record, nil
}
