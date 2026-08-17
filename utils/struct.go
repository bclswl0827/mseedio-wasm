package utils

import (
	"fmt"

	"github.com/bclswl0827/mseedio-wasm/wrapper/free"
	"github.com/bclswl0827/mseedio-wasm/wrapper/malloc"
	"github.com/tetratelabs/wazero/api"
)

type StructHandler struct {
	memoryObj  api.Memory
	mallocFunc *malloc.IMPL_malloc
	freeFunc   *free.IMPL_free
}

func NewStructHandler(memoryObj api.Memory, mallocFunc *malloc.IMPL_malloc, freeFunc *free.IMPL_free) StructHandler {
	return StructHandler{
		memoryObj:  memoryObj,
		mallocFunc: mallocFunc,
		freeFunc:   freeFunc,
	}
}

func (s *StructHandler) Pointer(size uint32) (uint32, error) {
	if s.mallocFunc == nil {
		return 0, fmt.Errorf("malloc function is nil")
	}
	if s.memoryObj == nil {
		return 0, fmt.Errorf("memory object is nil")
	}

	ptr, err := s.mallocFunc.Call(size)
	if err != nil {
		return 0, fmt.Errorf("failed to malloc memory: %w", err)
	}
	if ptr == 0 {
		return 0, fmt.Errorf("failed to allocate memory")
	}

	if ok := s.memoryObj.Write(ptr, make([]byte, size)); !ok {
		return 0, fmt.Errorf("failed to write memory")
	}

	return ptr, nil
}

func (s *StructHandler) Read(pointer uint32, size uint32) ([]byte, error) {
	if s.memoryObj == nil {
		return nil, fmt.Errorf("memory object is nil")
	}
	if pointer == 0 {
		return nil, fmt.Errorf("pointer is nil")
	}

	data, ok := s.memoryObj.Read(pointer, size)
	if !ok {
		return nil, fmt.Errorf("failed to read memory at 0x%x", pointer)
	}

	return data, nil
}

func (s *StructHandler) Free(pointer uint32) error {
	if s.freeFunc == nil {
		return fmt.Errorf("free function is nil")
	}
	if pointer == 0 {
		return fmt.Errorf("pointer is 0")
	}
	if err := s.freeFunc.Call(pointer); err != nil {
		return fmt.Errorf("failed to free struct: %w", err)
	}

	return nil
}
