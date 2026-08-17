package types

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/bclswl0827/mseedio-wasm/utils"
)

/*
Offset of record:        0
Offset of reclen:        4
Offset of swapflag:      8
Offset of sid:           9
Offset of formatversion: 73
Offset of flags:         74
Offset of starttime:     80
Offset of samprate:      88
Offset of encoding:      96
Offset of pubversion:    98
Offset of samplecnt:     104
Offset of crc:           112
Offset of extralength:   116
Offset of datalength:    120
Offset of extra:         124
Offset of datasamples:   128
Offset of datasize:      136
Offset of numsamples:    144
Offset of sampletype:    152
sizeof(MS3Record):       160
*/

const (
	LM_SIDLEN      = 64
	SizeOfM3Record = 160 // sizeof(M3Record) in C
)

type M3Record struct {
	Record   []byte
	RecLen   int32
	SwapFlag uint8

	Sid           [LM_SIDLEN]byte
	FormatVersion uint8
	Flags         uint8
	StartTime     int64
	SampleRate    float64
	Encoding      int16
	PubVersion    uint8
	SampleCnt     int64
	CRC           uint32
	ExtraLength   uint16
	DataLength    uint32
	Extra         []byte

	DataSamples []byte
	DataSize    uint64
	NumSamples  int64
	SampleType  byte
}

func (m *M3Record) FromBytes(structHandler utils.StructHandler, data []byte) error {
	if len(data) < SizeOfM3Record {
		return fmt.Errorf("buffer too small for M3Record")
	}

	recordAddr := binary.LittleEndian.Uint32(data[0:4]) // cast later if needed
	m.RecLen = int32(binary.LittleEndian.Uint32(data[4:8]))

	if recordAddr != 0 {
		recordBytes, err := structHandler.Read(recordAddr, uint32(m.RecLen))
		if err != nil {
			return err
		}
		m.Record = append([]byte(nil), recordBytes...)
	}

	m.SwapFlag = data[8]
	copy(m.Sid[:], data[9:73])

	m.FormatVersion = data[73]
	m.Flags = data[74]

	m.StartTime = int64(binary.LittleEndian.Uint64(data[80:88]))
	m.SampleRate = math.Float64frombits(binary.LittleEndian.Uint64(data[88:96]))
	m.Encoding = int16(binary.LittleEndian.Uint16(data[96:98]))
	m.PubVersion = data[98]

	m.SampleCnt = int64(binary.LittleEndian.Uint64(data[104:112]))
	m.CRC = binary.LittleEndian.Uint32(data[112:116])
	m.ExtraLength = binary.LittleEndian.Uint16(data[116:118])
	m.DataLength = binary.LittleEndian.Uint32(data[120:124])

	extraAddr := binary.LittleEndian.Uint32(data[124:128])
	if extraAddr != 0 {
		extraBytes, err := structHandler.Read(extraAddr, uint32(m.ExtraLength))
		if err != nil {
			return err
		}
		m.Extra = append([]byte(nil), extraBytes...)
	}

	m.DataSize = binary.LittleEndian.Uint64(data[136:144])
	m.NumSamples = int64(binary.LittleEndian.Uint64(data[144:152]))
	m.SampleType = data[152]

	dataSamplesAddr := binary.LittleEndian.Uint32(data[128:132])
	if dataSamplesAddr != 0 {
		if m.DataSize > math.MaxUint32 {
			return fmt.Errorf("data sample buffer too large: %d", m.DataSize)
		}
		dataSamplesBytes, err := structHandler.Read(dataSamplesAddr, uint32(m.DataSize))
		if err != nil {
			return err
		}
		m.DataSamples = append([]byte(nil), dataSamplesBytes...)
	}

	return nil
}
