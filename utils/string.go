package utils

import (
	"fmt"

	"github.com/bclswl0827/mseedio-wasm/wrapper/free"
	"github.com/bclswl0827/mseedio-wasm/wrapper/malloc"
	"github.com/tetratelabs/wazero/api"
)

type StringHandler struct {
	memoryObj  api.Memory
	mallocFunc *malloc.IMPL_malloc
	freeFunc   *free.IMPL_free
}

func NewStringHandler(memoryObj api.Memory, mallocFunc *malloc.IMPL_malloc, freeFunc *free.IMPL_free) StringHandler {
	return StringHandler{
		memoryObj:  memoryObj,
		mallocFunc: mallocFunc,
		freeFunc:   freeFunc,
	}
}

func (s *StringHandler) Pointer(str string) (uint32, error) {
	if s.mallocFunc == nil {
		return 0, fmt.Errorf("malloc function is nil")
	}
	if s.memoryObj == nil {
		return 0, fmt.Errorf("memory object is nil")
	}
	if str == "" {
		return 0, fmt.Errorf("string is empty")
	}

	strBuf := make([]byte, len(str)+1)
	copy(strBuf, []byte(str))
	strBuf[len(str)] = '\000'

	ptr, err := s.mallocFunc.Call(uint32(len(strBuf)))
	if err != nil {
		return 0, fmt.Errorf("failed to allocate memory: %w", err)
	}
	if ptr == 0 {
		return 0, fmt.Errorf("failed to allocate memory")
	}

	if ok := s.memoryObj.Write(ptr, strBuf); !ok {
		return 0, fmt.Errorf("failed to write memory")
	}

	return ptr, nil
}

func (s *StringHandler) Read(pointer uint32) (string, error) {
	if s.memoryObj == nil {
		return "", fmt.Errorf("memory object is nil")
	}
	if pointer == 0 {
		return "", fmt.Errorf("pointer is nil")
	}

	var bytes []byte
	for i := pointer; ; i++ {
		b, ok := s.memoryObj.ReadByte(i)
		if !ok {
			return "", fmt.Errorf("failed to read byte at 0x%x", i)
		}
		if b == '\000' {
			break
		}
		bytes = append(bytes, b)
	}

	return string(bytes), nil
}

func (s *StringHandler) Free(pointer uint32) error {
	if s.freeFunc == nil {
		return fmt.Errorf("free function is nil")
	}
	if pointer == 0 {
		return fmt.Errorf("pointer is nil")
	}
	if err := s.freeFunc.Call(pointer); err != nil {
		return fmt.Errorf("failed to free string: %w", err)
	}

	return nil
}
