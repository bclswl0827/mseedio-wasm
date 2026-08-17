package main

import (
	"fmt"
	"os"

	"github.com/bclswl0827/mseedio-wasm"
)

func main() {
	paths := os.Args[1:]
	if len(paths) == 0 {
		panic("miniseed path not specified")
	}

	library := mseedio.New()
	if err := library.Init(nil, nil); err != nil {
		panic(err)
	}
	defer library.Close()

	for _, path := range paths {
		if err := parseFile(library, path); err != nil {
			panic(err)
		}
	}
}

func parseFile(library *mseedio.MiniSEED, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	records, err := library.Parse(data)
	if err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}

	fmt.Printf("%s: %d bytes, %d records\n", path, len(data), len(records))
	for index, record := range records {
		start, err := library.FormatTime(record.StartTime, mseedio.TimeISOZ, mseedio.SubsecondsNanoOrMicroIfNonzero)
		if err != nil {
			return fmt.Errorf("format record %d start time: %w", index, err)
		}
		fmt.Printf("  %d: sid=%s version=%d start=%s rate=%g samples=%d type=%c\n",
			index, record.SID, record.FormatVersion, start, record.SampleRate,
			record.Samples.Len(), record.Samples.Type,
		)
	}
	return nil
}
