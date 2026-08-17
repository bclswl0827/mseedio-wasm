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

	sid, err := library.NSLCToSID(mseedio.NSLC{
		Network: "TW", Station: "TPUB", Location: "00", Channel: "BHZ",
	})
	if err != nil {
		panic(err)
	}
	components, err := library.SIDToNSLC(sid)
	if err != nil {
		panic(err)
	}
	fmt.Println(sid)
	fmt.Printf("%+v\n", components)
}
