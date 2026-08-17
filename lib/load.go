package lib

import (
	_ "embed"
)

//go:generate ./build_libmseed_wasm.sh

//go:embed libmseed.wasm
var wasmBytes []byte

func LoadBytes() []byte {
	return wasmBytes
}
