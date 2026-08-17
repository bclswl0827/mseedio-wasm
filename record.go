package mseedio

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math"
)

const (
	flagUnpackData  uint32 = 0x0001
	flagValidateCRC uint32 = 0x0004
	flagAtEndOfFile uint32 = 0x0010
	flagFlushData   uint32 = 0x0040
	flagPackV2      uint32 = 0x0080

	ms3RecordSize = uint32(160)
	maxSIDLength  = 64
)

// Encoding is a miniSEED payload encoding code.
type Encoding int16

const (
	EncodingAuto        Encoding = -1
	EncodingText        Encoding = 0
	EncodingInt16       Encoding = 1
	EncodingInt32       Encoding = 3
	EncodingFloat32     Encoding = 4
	EncodingFloat64     Encoding = 5
	EncodingSteim1      Encoding = 10
	EncodingSteim2      Encoding = 11
	EncodingGeoscope24  Encoding = 12
	EncodingGeoscope163 Encoding = 13
	EncodingGeoscope164 Encoding = 14
	EncodingCDSN        Encoding = 16
	EncodingSRO         Encoding = 30
	EncodingDWWSSN      Encoding = 32
)

// SampleType identifies the unpacked in-memory sample representation.
type SampleType byte

const (
	SampleTypeText    SampleType = 't'
	SampleTypeInt32   SampleType = 'i'
	SampleTypeFloat32 SampleType = 'f'
	SampleTypeFloat64 SampleType = 'd'
)

// Samples stores exactly one populated sample slice. Keeping the concrete
// representation avoids precision loss when integer records are decoded.
type Samples struct {
	Type    SampleType
	Text    []byte
	Int32   []int32
	Float32 []float32
	Float64 []float64
}

func (s Samples) Len() int {
	switch s.Type {
	case SampleTypeText:
		return len(s.Text)
	case SampleTypeInt32:
		return len(s.Int32)
	case SampleTypeFloat32:
		return len(s.Float32)
	case SampleTypeFloat64:
		return len(s.Float64)
	default:
		return 0
	}
}

// Record is an owned Go copy of an MS3Record. It remains valid after the next
// libmseed call and after MiniSEED.Close.
type Record struct {
	Raw             []byte
	SID             string
	FormatVersion   uint8
	Flags           uint8
	StartTime       NSTime
	EndTime         NSTime
	SampleRate      float64
	Encoding        Encoding
	Publication     uint8
	SampleCount     int64
	CRC             uint32
	ExtraHeaders    json.RawMessage
	EncodedDataSize uint32
	Samples         Samples
}

type ParseOptions struct {
	ValidateCRC bool
	Verbose     int8
}

// Parse parses every concatenated miniSEED record and unpacks its samples.
func (m *MiniSEED) Parse(data []byte) ([]Record, error) {
	return m.ParseWithOptions(data, ParseOptions{})
}

// ParseWithOptions parses every concatenated miniSEED record in data.
func (m *MiniSEED) ParseWithOptions(data []byte, options ParseOptions) ([]Record, error) {
	if err := m.ready(); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("mseedio: parse: empty input")
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	inputPtr, err := m.allocBytesLocked(data)
	if err != nil {
		return nil, err
	}
	defer m.freeLocked(inputPtr)

	flags := flagUnpackData | flagAtEndOfFile
	if options.ValidateCRC {
		flags |= flagValidateCRC
	}

	records := make([]Record, 0, 1)
	for offset := 0; offset < len(data); {
		record, consumed, parseErr := m.parseRecordLocked(
			inputPtr+uint32(offset), data[offset:], flags, options.Verbose,
		)
		if parseErr != nil {
			return nil, fmt.Errorf("record at byte %d: %w", offset, parseErr)
		}
		if consumed <= 0 || consumed > len(data)-offset {
			return nil, fmt.Errorf("mseedio: invalid consumed record length %d at byte %d", consumed, offset)
		}
		records = append(records, record)
		offset += consumed
	}
	return records, nil
}

// ParseRecord parses the first miniSEED record in data and unpacks its samples.
func (m *MiniSEED) ParseRecord(data []byte) (Record, error) {
	return m.ParseRecordWithOptions(data, ParseOptions{})
}

// ParseRecordWithOptions parses the first miniSEED record with explicit options.
func (m *MiniSEED) ParseRecordWithOptions(data []byte, options ParseOptions) (Record, error) {
	if err := m.ready(); err != nil {
		return Record{}, err
	}
	if len(data) == 0 {
		return Record{}, fmt.Errorf("mseedio: parse record: empty input")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	inputPtr, err := m.allocBytesLocked(data)
	if err != nil {
		return Record{}, err
	}
	defer m.freeLocked(inputPtr)

	flags := uint32(flagUnpackData | flagAtEndOfFile)
	if options.ValidateCRC {
		flags |= flagValidateCRC
	}
	record, _, err := m.parseRecordLocked(inputPtr, data, flags, options.Verbose)
	return record, err
}

// Detect returns the first record's byte length and format version without
// fully parsing or unpacking it.
func (m *MiniSEED) Detect(data []byte) (recordLength int64, formatVersion uint8, err error) {
	if readyErr := m.ready(); readyErr != nil {
		return 0, 0, readyErr
	}
	if len(data) == 0 {
		return 0, 0, fmt.Errorf("mseedio: detect: empty input")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	inputPtr, err := m.allocBytesLocked(data)
	if err != nil {
		return 0, 0, err
	}
	defer m.freeLocked(inputPtr)
	versionPtr, err := m.allocLocked(1)
	if err != nil {
		return 0, 0, err
	}
	defer m.freeLocked(versionPtr)

	results, err := m.callLocked(FnMS3Detect, uint64(inputPtr), uint64(len(data)), uint64(versionPtr))
	if err != nil {
		return 0, 0, err
	}
	version, ok := m.memory.ReadByte(versionPtr)
	if !ok {
		return 0, 0, fmt.Errorf("mseedio: read detected format version")
	}
	length := int64(results[0])
	if length < 0 {
		return length, version, m.libraryErrorLocked("detect record", NotMiniSEED)
	}
	return length, version, nil
}

func (m *MiniSEED) parseRecordLocked(inputPtr uint32, raw []byte, flags uint32, verbose int8) (Record, int, error) {
	recordPointerPointer, err := m.allocLocked(4)
	if err != nil {
		return Record{}, 0, err
	}
	defer m.freeLocked(recordPointerPointer)

	results, err := m.callLocked(
		FnMSR3Parse,
		uint64(inputPtr), uint64(len(raw)), uint64(recordPointerPointer), uint64(flags), uint64(uint8(verbose)),
	)
	if err != nil {
		return Record{}, 0, err
	}
	code := int32(results[0])
	if code > 0 {
		return Record{}, 0, &NeedMoreDataError{Needed: code}
	}
	if code < 0 {
		return Record{}, 0, m.libraryErrorLocked("parse record", code)
	}

	recordPointerBytes, err := m.readLocked(recordPointerPointer, 4)
	if err != nil {
		return Record{}, 0, err
	}
	recordPointer := binary.LittleEndian.Uint32(recordPointerBytes)
	if recordPointer == 0 {
		return Record{}, 0, fmt.Errorf("mseedio: parse record: libmseed returned NULL")
	}
	defer m.callLocked(FnMSR3Free, uint64(recordPointerPointer))

	record, err := m.copyRecordLocked(recordPointer, raw)
	if err != nil {
		return Record{}, 0, err
	}
	return record, len(record.Raw), nil
}

func (m *MiniSEED) copyRecordLocked(recordPointer uint32, raw []byte) (Record, error) {
	data, err := m.readLocked(recordPointer, ms3RecordSize)
	if err != nil {
		return Record{}, err
	}

	recordLength := int(int32(binary.LittleEndian.Uint32(data[4:8])))
	if recordLength <= 0 || recordLength > len(raw) {
		return Record{}, fmt.Errorf("mseedio: parsed record length %d exceeds %d-byte input", recordLength, len(raw))
	}
	sidBytes := data[9:73]
	if index := bytes.IndexByte(sidBytes, 0); index >= 0 {
		sidBytes = sidBytes[:index]
	}

	record := Record{
		Raw:             append([]byte(nil), raw[:recordLength]...),
		SID:             string(sidBytes),
		FormatVersion:   data[73],
		Flags:           data[74],
		StartTime:       NSTime(int64(binary.LittleEndian.Uint64(data[80:88]))),
		SampleRate:      math.Float64frombits(binary.LittleEndian.Uint64(data[88:96])),
		Encoding:        Encoding(int16(binary.LittleEndian.Uint16(data[96:98]))),
		Publication:     data[98],
		SampleCount:     int64(binary.LittleEndian.Uint64(data[104:112])),
		CRC:             binary.LittleEndian.Uint32(data[112:116]),
		EncodedDataSize: binary.LittleEndian.Uint32(data[120:124]),
	}

	extraLength := uint32(binary.LittleEndian.Uint16(data[116:118]))
	extraPointer := binary.LittleEndian.Uint32(data[124:128])
	if extraLength > 0 && extraPointer != 0 {
		extra, readErr := m.readLocked(extraPointer, extraLength)
		if readErr != nil {
			return Record{}, readErr
		}
		record.ExtraHeaders = json.RawMessage(extra)
	}

	results, callErr := m.callLocked(FnMSR3EndTime, uint64(recordPointer))
	if callErr != nil {
		return Record{}, callErr
	}
	record.EndTime = NSTime(int64(results[0]))

	samplePointer := binary.LittleEndian.Uint32(data[128:132])
	sampleBufferSize := binary.LittleEndian.Uint64(data[136:144])
	sampleCount := int64(binary.LittleEndian.Uint64(data[144:152]))
	sampleType := SampleType(data[152])
	if samplePointer != 0 && sampleCount > 0 {
		samples, decodeErr := m.copySamplesLocked(samplePointer, sampleBufferSize, sampleCount, sampleType)
		if decodeErr != nil {
			return Record{}, decodeErr
		}
		record.Samples = samples
	}
	return record, nil
}

func (m *MiniSEED) copySamplesLocked(pointer uint32, bufferSize uint64, count int64, sampleType SampleType) (Samples, error) {
	var sampleSize uint64
	switch sampleType {
	case SampleTypeText:
		sampleSize = 1
	case SampleTypeInt32, SampleTypeFloat32:
		sampleSize = 4
	case SampleTypeFloat64:
		sampleSize = 8
	default:
		return Samples{}, fmt.Errorf("mseedio: unsupported unpacked sample type %q", byte(sampleType))
	}
	if count < 0 || uint64(count) > math.MaxUint32/sampleSize {
		return Samples{}, fmt.Errorf("mseedio: invalid unpacked sample count %d", count)
	}
	needed := uint64(count) * sampleSize
	if needed > bufferSize {
		return Samples{}, fmt.Errorf("mseedio: sample data requires %d bytes, buffer has %d", needed, bufferSize)
	}
	data, err := m.readLocked(pointer, uint32(needed))
	if err != nil {
		return Samples{}, err
	}
	samples := Samples{Type: sampleType}
	switch sampleType {
	case SampleTypeText:
		samples.Text = data
	case SampleTypeInt32:
		samples.Int32 = make([]int32, count)
		for index := range samples.Int32 {
			samples.Int32[index] = int32(binary.LittleEndian.Uint32(data[index*4:]))
		}
	case SampleTypeFloat32:
		samples.Float32 = make([]float32, count)
		for index := range samples.Float32 {
			samples.Float32[index] = math.Float32frombits(binary.LittleEndian.Uint32(data[index*4:]))
		}
	case SampleTypeFloat64:
		samples.Float64 = make([]float64, count)
		for index := range samples.Float64 {
			samples.Float64[index] = math.Float64frombits(binary.LittleEndian.Uint64(data[index*8:]))
		}
	}
	return samples, nil
}

type PackOptions struct {
	FormatVersion uint8
	Encoding      Encoding
	RecordLength  int32
	Verbose       int8
}

// Pack packs a record and concatenates all generated miniSEED records.
func (m *MiniSEED) Pack(record Record, options *PackOptions) ([]byte, error) {
	records, err := m.PackRecords(record, options)
	if err != nil {
		return nil, err
	}
	return bytes.Join(records, nil), nil
}

// PackTo packs a record and writes all generated miniSEED records to dst.
func (m *MiniSEED) PackTo(record Record, options *PackOptions, dst io.Writer) error {
	if dst == nil {
		return fmt.Errorf("mseedio: pack: destination writer is nil")
	}
	records, err := m.PackRecords(record, options)
	if err != nil {
		return err
	}
	for _, data := range records {
		written, err := dst.Write(data)
		if err != nil {
			return fmt.Errorf("mseedio: pack: write output: %w", err)
		}
		if written != len(data) {
			return io.ErrShortWrite
		}
	}
	return nil
}

// PackRecords packs a record and returns each generated miniSEED record separately.
func (m *MiniSEED) PackRecords(record Record, options *PackOptions) ([][]byte, error) {
	if err := m.ready(); err != nil {
		return nil, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	resolved := PackOptions{FormatVersion: 3, Encoding: EncodingAuto, RecordLength: -1}
	if options != nil {
		resolved = *options
		if resolved.FormatVersion == 0 {
			resolved.FormatVersion = 3
		}
		if resolved.RecordLength == 0 {
			resolved.RecordLength = -1
		}
	}
	if resolved.FormatVersion != 2 && resolved.FormatVersion != 3 {
		return nil, fmt.Errorf("mseedio: pack: format version must be 2 or 3")
	}

	recordPointer, allocations, err := m.marshalRecordLocked(record, resolved)
	if err != nil {
		return nil, err
	}
	defer func() {
		for index := len(allocations) - 1; index >= 0; index-- {
			_ = m.freeLocked(allocations[index])
		}
	}()

	flags := flagFlushData
	if resolved.FormatVersion == 2 {
		flags |= flagPackV2
	}
	results, err := m.callLocked(FnMSR3PackInit, uint64(recordPointer), uint64(flags), uint64(uint8(resolved.Verbose)))
	if err != nil {
		return nil, err
	}
	packer := uint32(results[0])
	if packer == 0 {
		return nil, fmt.Errorf("mseedio: initialize record packer: libmseed returned NULL")
	}
	packerPointer, err := m.allocLocked(4)
	if err != nil {
		return nil, err
	}
	defer m.freeLocked(packerPointer)
	if !m.memory.WriteUint32Le(packerPointer, packer) {
		return nil, fmt.Errorf("mseedio: store record packer pointer")
	}
	defer m.callLocked(FnMSR3PackFree, uint64(packerPointer), 0)

	recordOutputPointer, err := m.allocLocked(4)
	if err != nil {
		return nil, err
	}
	defer m.freeLocked(recordOutputPointer)
	recordLengthPointer, err := m.allocLocked(4)
	if err != nil {
		return nil, err
	}
	defer m.freeLocked(recordLengthPointer)

	var packed [][]byte
	for {
		results, callErr := m.callLocked(FnMSR3PackNext, uint64(packer), uint64(recordOutputPointer), uint64(recordLengthPointer))
		if callErr != nil {
			return nil, callErr
		}
		status := int32(results[0])
		if status == 0 {
			break
		}
		if status < 0 {
			return nil, fmt.Errorf("mseedio: generate next packed record")
		}
		outputPointer, ok := m.memory.ReadUint32Le(recordOutputPointer)
		if !ok {
			return nil, fmt.Errorf("mseedio: read packed record pointer")
		}
		lengthValue, ok := m.memory.ReadUint32Le(recordLengthPointer)
		if !ok || int32(lengthValue) <= 0 {
			return nil, fmt.Errorf("mseedio: read packed record length")
		}
		output, readErr := m.readLocked(outputPointer, lengthValue)
		if readErr != nil {
			return nil, readErr
		}
		packed = append(packed, output)
	}
	return packed, nil
}

func (m *MiniSEED) marshalRecordLocked(record Record, options PackOptions) (uint32, []uint32, error) {
	if record.SID == "" || len(record.SID) >= maxSIDLength {
		return 0, nil, fmt.Errorf("mseedio: pack: SID length must be between 1 and %d bytes", maxSIDLength-1)
	}
	if len(record.ExtraHeaders) > math.MaxUint16 {
		return 0, nil, fmt.Errorf("mseedio: pack: extra headers exceed %d bytes", math.MaxUint16)
	}

	sampleData, sampleType, sampleCount, inferredEncoding, err := encodeSamples(record.Samples)
	if err != nil {
		return 0, nil, err
	}
	encoding := options.Encoding
	if encoding == EncodingAuto || (encoding == EncodingText && sampleType != SampleTypeText) {
		encoding = inferredEncoding
	}

	allocations := make([]uint32, 0, 3)
	cleanup := func() {
		for index := len(allocations) - 1; index >= 0; index-- {
			_ = m.freeLocked(allocations[index])
		}
	}

	extraPointer := uint32(0)
	if len(record.ExtraHeaders) > 0 {
		extraPointer, err = m.allocBytesLocked(record.ExtraHeaders)
		if err != nil {
			return 0, nil, err
		}
		allocations = append(allocations, extraPointer)
	}
	samplePointer := uint32(0)
	if len(sampleData) > 0 {
		samplePointer, err = m.allocBytesLocked(sampleData)
		if err != nil {
			cleanup()
			return 0, nil, err
		}
		allocations = append(allocations, samplePointer)
	}
	recordPointer, err := m.allocLocked(ms3RecordSize)
	if err != nil {
		cleanup()
		return 0, nil, err
	}
	allocations = append(allocations, recordPointer)

	data := make([]byte, ms3RecordSize)
	binary.LittleEndian.PutUint32(data[4:8], uint32(options.RecordLength))
	copy(data[9:73], record.SID)
	data[73] = options.FormatVersion
	data[74] = record.Flags
	binary.LittleEndian.PutUint64(data[80:88], uint64(record.StartTime))
	binary.LittleEndian.PutUint64(data[88:96], math.Float64bits(record.SampleRate))
	binary.LittleEndian.PutUint16(data[96:98], uint16(encoding))
	data[98] = record.Publication
	binary.LittleEndian.PutUint64(data[104:112], uint64(sampleCount))
	binary.LittleEndian.PutUint16(data[116:118], uint16(len(record.ExtraHeaders)))
	binary.LittleEndian.PutUint32(data[124:128], extraPointer)
	binary.LittleEndian.PutUint32(data[128:132], samplePointer)
	binary.LittleEndian.PutUint64(data[136:144], uint64(len(sampleData)))
	binary.LittleEndian.PutUint64(data[144:152], uint64(sampleCount))
	data[152] = byte(sampleType)
	if !m.memory.Write(recordPointer, data) {
		cleanup()
		return 0, nil, fmt.Errorf("mseedio: write MS3Record at %#x", recordPointer)
	}
	return recordPointer, allocations, nil
}

func encodeSamples(samples Samples) ([]byte, SampleType, int64, Encoding, error) {
	count := samples.Len()
	if count == 0 {
		return nil, 0, 0, EncodingText, nil
	}
	var data []byte
	switch samples.Type {
	case SampleTypeText:
		data = append([]byte(nil), samples.Text...)
		return data, samples.Type, int64(count), EncodingText, nil
	case SampleTypeInt32:
		data = make([]byte, count*4)
		for index, value := range samples.Int32 {
			binary.LittleEndian.PutUint32(data[index*4:], uint32(value))
		}
		return data, samples.Type, int64(count), EncodingSteim2, nil
	case SampleTypeFloat32:
		data = make([]byte, count*4)
		for index, value := range samples.Float32 {
			binary.LittleEndian.PutUint32(data[index*4:], math.Float32bits(value))
		}
		return data, samples.Type, int64(count), EncodingFloat32, nil
	case SampleTypeFloat64:
		data = make([]byte, count*8)
		for index, value := range samples.Float64 {
			binary.LittleEndian.PutUint64(data[index*8:], math.Float64bits(value))
		}
		return data, samples.Type, int64(count), EncodingFloat64, nil
	default:
		return nil, 0, 0, EncodingAuto, fmt.Errorf("mseedio: pack: unknown sample type %q", byte(samples.Type))
	}
}
