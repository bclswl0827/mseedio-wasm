# mseedio-wasm

CGO-free Go bindings for [libmseed](https://github.com/EarthScope/libmseed).

## Build

The upstream libmseed source is tracked as a Git submodule. After cloning this repository, initialize it with:

```sh
$ git submodule update --init --recursive
```

Alternatively, simply cloning this repository with `--recurse-submodules`.

Next, build the embedded WASM module (requires `make` and Emscripten's `emcc`):

```sh
$ ./lib/build_libmseed_wasm.sh
```

## APIs

- `MiniSEED.Parse`, `ParseRecord`, `Detect`, `Pack` and `PackRecords` own all WASM memory and return ordinary Go values.
- `ParseTime`, `FormatTime`, `SIDToNSLC`, `NSLCToSID` and `NormalizeSID` cover common utility operations.
- `MiniSEED.Raw` exposes the exported libmseed functions listed in [lib/README.md](lib/README.md). Raw calls use the WASM ABI and are intended for advanced APIs such as selections, trace lists and extra headers.

```go
library := mseedio.New()
if err := library.Init(nil, nil); err != nil {
    return err
}
defer library.Close()

records, err := library.Parse(input)
```

The values returned by the high-level API are copies and remain valid after subsequent calls or `Close`. A single `MiniSEED` instance serializes access to its WASM module and can be shared by goroutines.

## Examples

Each example is an independent command:

- `go run ./example/parse` parses upstream libmseed MiniSEED 2 and 3 test data; optional file arguments override the defaults.
- `go run ./example/pack` writes `demo.mseed2` and `demo.mseed3`.
- `go run ./example/time`
- `go run ./example/sid`
- `go run ./example/raw`

## Raw ABI

`Raw.Call` accepts and returns wazero ABI values. Allocate guest pointers with `Raw.Alloc`, `Raw.AllocBytes` or `Raw.AllocString`, and release owned pointers with `Raw.Free`. Use `api.EncodeF32`, `api.EncodeF64`, `api.DecodeF32` and `api.DecodeF64` for floating-point values.

C-callback packing functions are intentionally not exported. Record packing uses `msr3_pack_init`, `msr3_pack_next` and `msr3_pack_free`; trace-list packing has an equivalent generator API.

Raw file-based functions run inside WASI and see only preopened paths. Mount a directory explicitly with `InitWithOptions` and a `wazero.FSConfig`; the buffer-based high-level parse/pack API does not need filesystem access.

## Third-party code

The WASM module is built directly from the unmodified upstream [EarthScope/libmseed](https://github.com/EarthScope/libmseed), licensed under the Apache License 2.0. See `lib/LICENSE`.

## License

This project is licensed under the MIT License.
