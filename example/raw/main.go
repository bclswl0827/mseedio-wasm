package main

import (
	"fmt"

	"github.com/bclswl0827/mseedio-wasm"
)

func main() {
	library := mseedio.New()
	if err := library.Init(nil, nil); err != nil {
		panic(err)
	}
	defer library.Close()

	results, err := library.Raw.Call(mseedio.FnMSSampleSize, uint64('i'))
	if err != nil {
		panic(err)
	}
	fmt.Printf("int32 sample size: %d bytes\n", results[0])
	fmt.Printf("available libmseed function exports: %d\n", len(library.Raw.Names()))
}
