package utils

import (
	"encoding/binary"
	"fmt"
	"unsafe"

	"github.com/bclswl0827/mseedio-wasm/wrapper/free"
	"github.com/bclswl0827/mseedio-wasm/wrapper/malloc"
	"github.com/tetratelabs/wazero/api"
)

type ValuePointerHandler struct {
	memoryObj  api.Memory
	mallocFunc *malloc.IMPL_malloc
	freeFunc   *free.IMPL_free
}

func NewValuePointerHandler(memoryObj api.Memory, mallocFunc *malloc.IMPL_malloc, freeFunc *free.IMPL_free) ValuePointerHandler {
	return ValuePointerHandler{
		memoryObj:  memoryObj,
		mallocFunc: mallocFunc,
		freeFunc:   freeFunc,
	}
}

func (v *ValuePointerHandler) Pointer(value any) (uint32, error) {
	if v.mallocFunc == nil {
		return 0, fmt.Errorf("malloc function is nil")
	}
	if v.memoryObj == nil {
		return 0, fmt.Errorf("memory object is nil")
	}

	size := uint32(binary.Size(value))
	ptr, err := v.mallocFunc.Call(size)
	if err != nil {
		return 0, fmt.Errorf("failed to allocate memory: %w", err)
	}
	if ptr == 0 {
		return 0, fmt.Errorf("failed to allocate memory")
	}

	buf := make([]byte, size)
	switch val := value.(type) {
	case uint8:
		buf[0] = val
	case uint16:
		binary.LittleEndian.PutUint16(buf, val)
	case int16:
		binary.LittleEndian.PutUint16(buf, uint16(val))
	case uint32:
		binary.LittleEndian.PutUint32(buf, val)
	case int32:
		binary.LittleEndian.PutUint32(buf, uint32(val))
	case uint64:
		binary.LittleEndian.PutUint64(buf, val)
	case int64:
		binary.LittleEndian.PutUint64(buf, uint64(val))
	case float32:
		binary.LittleEndian.PutUint32(buf, *(*uint32)(unsafe.Pointer(&val)))
	case float64:
		binary.LittleEndian.PutUint64(buf, *(*uint64)(unsafe.Pointer(&val)))
	default:
		return 0, fmt.Errorf("unsupported type: %T", val)
	}

	if ok := v.memoryObj.Write(ptr, buf); !ok {
		return 0, fmt.Errorf("failed to write memory")
	}

	return ptr, nil
}

func (v *ValuePointerHandler) Read(pointer uint32, result *any) error {
	if v.memoryObj == nil {
		return fmt.Errorf("memory object is nil")
	}
	if result == nil {
		return fmt.Errorf("result is nil")
	}
	if pointer == 0 {
		return fmt.Errorf("pointer is nil")
	}

	size := uint32(binary.Size(*result))
	data, ok := v.memoryObj.Read(pointer, size)
	if !ok {
		return fmt.Errorf("failed to read memory at 0x%x", pointer)
	}

	switch (*result).(type) {
	case uint8:
		*result = data[0]
		return nil
	case uint16:
		*result = binary.LittleEndian.Uint16(data)
		return nil
	case int16:
		*result = int16(binary.LittleEndian.Uint16(data))
		return nil
	case uint32:
		*result = binary.LittleEndian.Uint32(data)
		return nil
	case int32:
		*result = int32(binary.LittleEndian.Uint32(data))
		return nil
	case uint64:
		*result = binary.LittleEndian.Uint64(data)
		return nil
	case int64:
		*result = int64(binary.LittleEndian.Uint64(data))
		return nil
	case float32:
		bits := binary.LittleEndian.Uint32(data)
		*result = *(*float32)(unsafe.Pointer(&bits))
		return nil
	case float64:
		bits := binary.LittleEndian.Uint64(data)
		*result = *(*float64)(unsafe.Pointer(&bits))
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", result)
	}
}

func (v *ValuePointerHandler) Free(pointer uint32) error {
	if v.freeFunc == nil {
		return fmt.Errorf("free function is nil")
	}
	if pointer == 0 {
		return fmt.Errorf("pointer is 0")
	}
	if err := v.freeFunc.Call(pointer); err != nil {
		return fmt.Errorf("failed to free value: %w", err)
	}

	return nil
}
