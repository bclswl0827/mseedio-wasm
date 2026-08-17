package main

import (
	"bytes"
	"fmt"
	"math"
	"os"

	"github.com/bclswl0827/mseedio-wasm"
)

func main() {
	library := mseedio.New()
	if err := library.Init(nil, nil); err != nil {
		panic(err)
	}
	defer library.Close()

	start, err := library.ParseTime("2026-08-17T12:00:00Z")
	if err != nil {
		panic(err)
	}
	sid, err := library.NormalizeSID("AS.SHAKE.00.EHZ")
	if err != nil {
		panic(err)
	}

	const (
		sampleRate = 100.0
		duration   = 60.0
		frequency  = 1.0
		amplitude  = 1000.0
	)
	samples := make([]int32, int(sampleRate*duration))
	for index := range samples {
		time := float64(index) / sampleRate
		samples[index] = int32(math.Round(amplitude * math.Sin(2*math.Pi*frequency*time)))
	}

	record := mseedio.Record{
		SID:        sid,
		StartTime:  start,
		SampleRate: sampleRate,
		Samples: mseedio.Samples{
			Type:  mseedio.SampleTypeInt32,
			Int32: samples,
		},
	}
	outputs := []struct {
		version uint8
		path    string
	}{
		{version: 2, path: "demo.mseed2"},
		{version: 3, path: "demo.mseed3"},
	}

	for _, output := range outputs {
		var buffer bytes.Buffer
		err := library.PackTo(record, &mseedio.PackOptions{
			FormatVersion: output.version,
			Encoding:      mseedio.EncodingSteim2,
		}, &buffer)
		if err != nil {
			panic(err)
		}
		if err := os.WriteFile(output.path, buffer.Bytes(), 0o644); err != nil {
			panic(err)
		}
		fmt.Printf("wrote %d bytes of miniSEED %d to %s\n", buffer.Len(), output.version, output.path)
	}
}
