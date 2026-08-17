package mseedio

import (
	"fmt"
	"math"
	"sort"

	"github.com/tetratelabs/wazero/api"
)

// FunctionName identifies a C function exported by libmseed.wasm.
type FunctionName string

// Function names retained by the Go/WASM binding export policy.
const (
	FnMalloc                     FunctionName = "malloc"
	FnFree                       FunctionName = "free"
	FnMS3MSFPInit                FunctionName = "ms3_msfp_init"
	FnMS3MSFPInitFD              FunctionName = "ms3_msfp_init_fd"
	FnMS3Detect                  FunctionName = "ms3_detect"
	FnMS3MatchSelect             FunctionName = "ms3_matchselect"
	FnMS3AddSelect               FunctionName = "ms3_addselect"
	FnMS3AddSelectComp           FunctionName = "ms3_addselect_comp"
	FnMS3ReadSelectionsFile      FunctionName = "ms3_readselectionsfile"
	FnMS3FreeSelections          FunctionName = "ms3_freeselections"
	FnMS3PrintSelections         FunctionName = "ms3_printselections"
	FnMS3ReadMSR                 FunctionName = "ms3_readmsr"
	FnMS3ReadMSRR                FunctionName = "ms3_readmsr_r"
	FnMS3ReadMSRSelection        FunctionName = "ms3_readmsr_selection"
	FnMS3ReadTraceList           FunctionName = "ms3_readtracelist"
	FnMS3ReadTraceListSelection  FunctionName = "ms3_readtracelist_selection"
	FnMS3ReadTraceListTimeWindow FunctionName = "ms3_readtracelist_timewin"

	FnMSR3Init           FunctionName = "msr3_init"
	FnMSR3Free           FunctionName = "msr3_free"
	FnMSR3Parse          FunctionName = "msr3_parse"
	FnMSR3EndTime        FunctionName = "msr3_endtime"
	FnMSR3NSPeriod       FunctionName = "msr3_nsperiod"
	FnMSR3SampleRateHz   FunctionName = "msr3_sampratehz"
	FnMSR3HostLatency    FunctionName = "msr3_host_latency"
	FnMSR3UnpackData     FunctionName = "msr3_unpack_data"
	FnMSR3DataBounds     FunctionName = "msr3_data_bounds"
	FnMSR3Duplicate      FunctionName = "msr3_duplicate"
	FnMSR3DuplicateExtra FunctionName = "msr3_duplicate_extra"
	FnMSR3ResizeBuffer   FunctionName = "msr3_resize_buffer"
	FnMSR3Print          FunctionName = "msr3_print"
	FnMSR3MatchSelect    FunctionName = "msr3_matchselect"
	FnMSR3PackInit       FunctionName = "msr3_pack_init"
	FnMSR3PackNext       FunctionName = "msr3_pack_next"
	FnMSR3PackFree       FunctionName = "msr3_pack_free"
	FnMSR3PackHeader2    FunctionName = "msr3_pack_header2"
	FnMSR3PackHeader3    FunctionName = "msr3_pack_header3"
	FnMSR3RepackMSeed2   FunctionName = "msr3_repack_mseed2"
	FnMSR3RepackMSeed3   FunctionName = "msr3_repack_mseed3"
	FnMSR3WriteMSeed     FunctionName = "msr3_writemseed"

	FnMSTL3Init                FunctionName = "mstl3_init"
	FnMSTL3Free                FunctionName = "mstl3_free"
	FnMSTL3FindID              FunctionName = "mstl3_findID"
	FnMSTL3AddMSR              FunctionName = "mstl3_addmsr"
	FnMSTL3AddMSRRecordPointer FunctionName = "mstl3_addmsr_recordptr"
	FnMSTL3ReadBuffer          FunctionName = "mstl3_readbuffer"
	FnMSTL3ReadBufferSelection FunctionName = "mstl3_readbuffer_selection"
	FnMSTL3UnpackRecordList    FunctionName = "mstl3_unpack_recordlist"
	FnMSTL3ConvertSamples      FunctionName = "mstl3_convertsamples"
	FnMSTL3ResizeBuffers       FunctionName = "mstl3_resize_buffers"
	FnMSTL3PackInit            FunctionName = "mstl3_pack_init"
	FnMSTL3PackNext            FunctionName = "mstl3_pack_next"
	FnMSTL3PackFree            FunctionName = "mstl3_pack_free"
	FnMSTL3PrintTraceList      FunctionName = "mstl3_printtracelist"
	FnMSTL3PrintSyncList       FunctionName = "mstl3_printsynclist"
	FnMSTL3PrintGapList        FunctionName = "mstl3_printgaplist"
	FnMSTL3WriteMSeed          FunctionName = "mstl3_writemseed"

	FnMSDecodeData         FunctionName = "ms_decode_data"
	FnMSParseRaw2          FunctionName = "ms_parse_raw2"
	FnMSParseRaw3          FunctionName = "ms_parse_raw3"
	FnMSCRC32C             FunctionName = "ms_crc32c"
	FnMSSampleSize         FunctionName = "ms_samplesize"
	FnMSEncodingSizeType   FunctionName = "ms_encoding_sizetype"
	FnMSEncodingString     FunctionName = "ms_encodingstr"
	FnMSErrorString        FunctionName = "ms_errorstr"
	FnMSSampleTime         FunctionName = "ms_sampletime"
	FnMSReadLeapSeconds    FunctionName = "ms_readleapseconds"
	FnMSReadLeapSecondFile FunctionName = "ms_readleapsecondfile"
	FnMSSidToNSLCN         FunctionName = "ms_sid2nslc_n"
	FnMSNSLCToSid          FunctionName = "ms_nslc2sid"
	FnMSSeedChannelToX     FunctionName = "ms_seedchan2xchan"
	FnMSXChannelToSeed     FunctionName = "ms_xchan2seedchan"

	FnMSNSTimeToTime        FunctionName = "ms_nstime2time"
	FnMSNSTimeToTimeStringN FunctionName = "ms_nstime2timestr_n"
	FnMSTimeToNSTime        FunctionName = "ms_time2nstime"
	FnMSTimeStringToNSTime  FunctionName = "ms_timestr2nstime"
	FnMSMDTimeToNSTime      FunctionName = "ms_mdtimestr2nstime"
	FnMSSeedTimeToNSTime    FunctionName = "ms_seedtimestr2nstime"
	FnMSDayOfYearToMD       FunctionName = "ms_doy2md"
	FnMSMDToDayOfYear       FunctionName = "ms_md2doy"

	FnMSEHGetPointerType     FunctionName = "mseh_get_ptr_type"
	FnMSEHGetPointer         FunctionName = "mseh_get_ptr_r"
	FnMSEHSetPointer         FunctionName = "mseh_set_ptr_r"
	FnMSEHSerialize          FunctionName = "mseh_serialize"
	FnMSEHFreeParseState     FunctionName = "mseh_free_parsestate"
	FnMSEHAddEventDetection  FunctionName = "mseh_add_event_detection_r"
	FnMSEHAddCalibration     FunctionName = "mseh_add_calibration_r"
	FnMSEHAddTimingException FunctionName = "mseh_add_timing_exception_r"
	FnMSEHAddRecenter        FunctionName = "mseh_add_recenter_r"
	FnMSEHReplace            FunctionName = "mseh_replace"
	FnMSEHPrint              FunctionName = "mseh_print"
)

var libmseedFunctionNames = []FunctionName{
	FnMalloc, FnFree,
	FnMS3MSFPInit, FnMS3MSFPInitFD, FnMS3Detect, FnMS3MatchSelect, FnMS3AddSelect,
	FnMS3AddSelectComp, FnMS3ReadSelectionsFile, FnMS3FreeSelections, FnMS3PrintSelections,
	FnMS3ReadMSR, FnMS3ReadMSRR, FnMS3ReadMSRSelection, FnMS3ReadTraceList,
	FnMS3ReadTraceListSelection, FnMS3ReadTraceListTimeWindow,
	FnMSR3Init, FnMSR3Free, FnMSR3Parse, FnMSR3EndTime, FnMSR3NSPeriod,
	FnMSR3SampleRateHz, FnMSR3HostLatency, FnMSR3UnpackData, FnMSR3DataBounds,
	FnMSR3Duplicate, FnMSR3DuplicateExtra, FnMSR3ResizeBuffer, FnMSR3Print,
	FnMSR3MatchSelect, FnMSR3PackInit, FnMSR3PackNext, FnMSR3PackFree,
	FnMSR3PackHeader2, FnMSR3PackHeader3, FnMSR3RepackMSeed2, FnMSR3RepackMSeed3,
	FnMSR3WriteMSeed, FnMSTL3Init, FnMSTL3Free, FnMSTL3FindID, FnMSTL3AddMSR,
	FnMSTL3AddMSRRecordPointer, FnMSTL3ReadBuffer, FnMSTL3ReadBufferSelection,
	FnMSTL3UnpackRecordList, FnMSTL3ConvertSamples, FnMSTL3ResizeBuffers,
	FnMSTL3PackInit, FnMSTL3PackNext, FnMSTL3PackFree,
	FnMSTL3PrintTraceList, FnMSTL3PrintSyncList, FnMSTL3PrintGapList,
	FnMSTL3WriteMSeed, FnMSDecodeData, FnMSParseRaw2, FnMSParseRaw3, FnMSCRC32C,
	FnMSSampleSize, FnMSEncodingSizeType, FnMSEncodingString, FnMSErrorString,
	FnMSSampleTime, FnMSReadLeapSeconds, FnMSReadLeapSecondFile, FnMSSidToNSLCN,
	FnMSNSLCToSid, FnMSSeedChannelToX, FnMSXChannelToSeed, FnMSNSTimeToTime,
	FnMSNSTimeToTimeStringN, FnMSTimeToNSTime, FnMSTimeStringToNSTime,
	FnMSMDTimeToNSTime, FnMSSeedTimeToNSTime, FnMSDayOfYearToMD,
	FnMSMDToDayOfYear, FnMSEHGetPointerType, FnMSEHGetPointer, FnMSEHSetPointer,
	FnMSEHSerialize, FnMSEHFreeParseState, FnMSEHAddEventDetection,
	FnMSEHAddCalibration, FnMSEHAddTimingException, FnMSEHAddRecenter,
	FnMSEHReplace, FnMSEHPrint,
}

// RawAPI is the CGO-free, one-to-one libmseed WASM layer. Parameters and
// results use the WebAssembly ABI: pointers/i32 values are uint64 values with
// the significant bits in the low 32 bits, while i64 and f64 use wazero's
// standard encoding. Use api.EncodeF64/api.DecodeF64 for floating-point ABI
// values.
type RawAPI struct {
	owner     *MiniSEED
	functions map[FunctionName]api.Function
}

func newRawAPI(owner *MiniSEED) (*RawAPI, error) {
	raw := &RawAPI{owner: owner, functions: make(map[FunctionName]api.Function)}
	for _, name := range libmseedFunctionNames {
		function := owner.module.ExportedFunction(string(name))
		if function == nil {
			return nil, fmt.Errorf("mseedio: expected function %q is missing from WASM", name)
		}
		raw.functions[name] = function
	}
	return raw, nil
}

// Names returns the available libmseed function exports in lexical order.
func (r *RawAPI) Names() []FunctionName {
	if r == nil {
		return nil
	}
	names := make([]FunctionName, 0, len(r.functions))
	for name := range r.functions {
		names = append(names, name)
	}
	sort.Slice(names, func(i, j int) bool { return names[i] < names[j] })
	return names
}

// Definition returns the wazero ABI metadata for an exported function.
func (r *RawAPI) Definition(name FunctionName) (api.FunctionDefinition, error) {
	if err := r.ready(); err != nil {
		return nil, err
	}
	function, ok := r.functions[name]
	if !ok {
		return nil, fmt.Errorf("mseedio: function %q is not exported", name)
	}
	return function.Definition(), nil
}

// Call invokes one exported function using raw WebAssembly ABI values.
func (r *RawAPI) Call(name FunctionName, params ...uint64) ([]uint64, error) {
	if err := r.ready(); err != nil {
		return nil, err
	}
	function, ok := r.functions[name]
	if !ok {
		return nil, fmt.Errorf("mseedio: function %q is not exported", name)
	}

	r.owner.mu.Lock()
	defer r.owner.mu.Unlock()
	results, err := function.Call(r.owner.ctx, params...)
	if err != nil {
		return nil, fmt.Errorf("mseedio: call %s: %w", name, err)
	}
	return results, nil
}

// Alloc allocates zeroed memory in the WASM address space.
func (r *RawAPI) Alloc(size uint32) (uint32, error) {
	if err := r.ready(); err != nil {
		return 0, err
	}
	if size == 0 {
		size = 1
	}
	r.owner.mu.Lock()
	defer r.owner.mu.Unlock()
	return r.owner.allocLocked(size)
}

// AllocBytes copies data into newly allocated WASM memory.
func (r *RawAPI) AllocBytes(data []byte) (uint32, error) {
	if err := r.ready(); err != nil {
		return 0, err
	}
	if uint64(len(data)) > math.MaxUint32 {
		return 0, fmt.Errorf("mseedio: byte slice exceeds the WASM32 address space")
	}
	size := uint32(len(data))
	if size == 0 {
		size = 1
	}
	r.owner.mu.Lock()
	defer r.owner.mu.Unlock()
	ptr, err := r.owner.allocLocked(size)
	if err != nil {
		return 0, err
	}
	if len(data) > 0 && !r.owner.memory.Write(ptr, data) {
		_ = r.owner.freeLocked(ptr)
		return 0, fmt.Errorf("mseedio: write WASM memory at %#x", ptr)
	}
	return ptr, nil
}

// AllocString copies a NUL-terminated string into newly allocated WASM memory.
func (r *RawAPI) AllocString(value string) (uint32, error) {
	data := append([]byte(value), 0)
	return r.AllocBytes(data)
}

// Free releases memory previously returned by Alloc, AllocBytes or AllocString.
func (r *RawAPI) Free(ptr uint32) error {
	if err := r.ready(); err != nil {
		return err
	}
	if ptr == 0 {
		return nil
	}
	r.owner.mu.Lock()
	defer r.owner.mu.Unlock()
	return r.owner.freeLocked(ptr)
}

// Read returns a copy of size bytes from WASM memory.
func (r *RawAPI) Read(ptr, size uint32) ([]byte, error) {
	if err := r.ready(); err != nil {
		return nil, err
	}
	r.owner.mu.Lock()
	defer r.owner.mu.Unlock()
	return r.owner.readLocked(ptr, size)
}

// Write copies data into WASM memory.
func (r *RawAPI) Write(ptr uint32, data []byte) error {
	if err := r.ready(); err != nil {
		return err
	}
	r.owner.mu.Lock()
	defer r.owner.mu.Unlock()
	if !r.owner.memory.Write(ptr, data) {
		return fmt.Errorf("mseedio: write %d bytes at %#x", len(data), ptr)
	}
	return nil
}

// ReadString reads a NUL-terminated string from WASM memory.
func (r *RawAPI) ReadString(ptr uint32) (string, error) {
	if err := r.ready(); err != nil {
		return "", err
	}
	r.owner.mu.Lock()
	defer r.owner.mu.Unlock()
	return r.owner.readStringLocked(ptr)
}

func (r *RawAPI) ready() error {
	if r == nil || r.owner == nil || !r.owner.initialized || r.owner.memory == nil {
		return fmt.Errorf("mseedio: runtime is not initialized")
	}
	return nil
}

func (m *MiniSEED) callLocked(name FunctionName, params ...uint64) ([]uint64, error) {
	function, ok := m.Raw.functions[name]
	if !ok {
		return nil, fmt.Errorf("mseedio: function %q is not exported", name)
	}
	results, err := function.Call(m.ctx, params...)
	if err != nil {
		return nil, fmt.Errorf("mseedio: call %s: %w", name, err)
	}
	return results, nil
}

func (m *MiniSEED) allocLocked(size uint32) (uint32, error) {
	results, err := m.callLocked(FnMalloc, uint64(size))
	if err != nil {
		return 0, err
	}
	ptr := uint32(results[0])
	if ptr == 0 {
		return 0, fmt.Errorf("mseedio: allocate %d bytes: out of memory", size)
	}
	if !m.memory.Write(ptr, make([]byte, size)) {
		_ = m.freeLocked(ptr)
		return 0, fmt.Errorf("mseedio: zero %d bytes at %#x", size, ptr)
	}
	return ptr, nil
}

func (m *MiniSEED) allocBytesLocked(data []byte) (uint32, error) {
	if uint64(len(data)) > math.MaxUint32 {
		return 0, fmt.Errorf("mseedio: byte slice exceeds the WASM32 address space")
	}
	size := uint32(len(data))
	if size == 0 {
		size = 1
	}
	ptr, err := m.allocLocked(size)
	if err != nil {
		return 0, err
	}
	if len(data) > 0 && !m.memory.Write(ptr, data) {
		_ = m.freeLocked(ptr)
		return 0, fmt.Errorf("mseedio: write WASM memory at %#x", ptr)
	}
	return ptr, nil
}

func (m *MiniSEED) freeLocked(ptr uint32) error {
	if ptr == 0 {
		return nil
	}
	_, err := m.callLocked(FnFree, uint64(ptr))
	return err
}

func (m *MiniSEED) readLocked(ptr, size uint32) ([]byte, error) {
	data, ok := m.memory.Read(ptr, size)
	if !ok {
		return nil, fmt.Errorf("mseedio: read %d bytes at %#x", size, ptr)
	}
	return append([]byte(nil), data...), nil
}

func (m *MiniSEED) readStringLocked(ptr uint32) (string, error) {
	if ptr == 0 {
		return "", fmt.Errorf("mseedio: read string from NULL")
	}
	data := make([]byte, 0, 64)
	for offset := uint32(0); ; offset++ {
		value, ok := m.memory.ReadByte(ptr + offset)
		if !ok {
			return "", fmt.Errorf("mseedio: read string byte at %#x", ptr+offset)
		}
		if value == 0 {
			return string(data), nil
		}
		data = append(data, value)
	}
}
