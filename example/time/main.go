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

	value, err := library.ParseTime("2026,229,12:34:56.123456789")
	if err != nil {
		panic(err)
	}
	formatted, err := library.FormatTime(value, mseedio.TimeISOZ, mseedio.SubsecondsNano)
	if err != nil {
		panic(err)
	}
	fmt.Println(formatted)
	fmt.Println(value.Time())
}
