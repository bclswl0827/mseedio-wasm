package mseedio

import (
	"fmt"
	"time"
)

const (
	NSTimeError NSTime = -2145916800000000000
	NSTimeUnset NSTime = -2145916799999999999
)

// NSTime is libmseed's nanoseconds-since-Unix-epoch time representation.
type NSTime int64

// TimeFormat controls FormatTime output.
type TimeFormat int32

const (
	TimeISO TimeFormat = iota
	TimeISOZ
	TimeISODayOfYear
	TimeISODayOfYearZ
	TimeISOSpace
	TimeISOSpaceZ
	TimeSEEDOrdinal
	TimeUnixEpoch
	TimeNanosecondEpoch
)

// SubsecondFormat controls the fractional-second part of FormatTime output.
type SubsecondFormat int32

const (
	SubsecondsNone SubsecondFormat = iota
	SubsecondsMicro
	SubsecondsNano
	SubsecondsMicroIfNonzero
	SubsecondsNanoIfNonzero
	SubsecondsNanoOrMicro
	SubsecondsNanoOrMicroIfNonzero
)

// NSTimeFromTime converts a time.Time to libmseed's nanoseconds-since-Unix-epoch value.
func NSTimeFromTime(value time.Time) NSTime {
	return NSTime(value.Unix()*int64(time.Second) + int64(value.Nanosecond()))
}

func (n NSTime) Time() time.Time {
	seconds := int64(n) / int64(time.Second)
	nanoseconds := int64(n) % int64(time.Second)
	return time.Unix(seconds, nanoseconds).UTC()
}

// ParseTime parses ISO, month/day and SEED ordinal formats accepted by libmseed.
func (m *MiniSEED) ParseTime(value string) (NSTime, error) {
	if err := m.ready(); err != nil {
		return 0, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	ptr, err := m.allocBytesLocked(append([]byte(value), 0))
	if err != nil {
		return 0, err
	}
	defer m.freeLocked(ptr)

	results, err := m.callLocked(FnMSTimeStringToNSTime, uint64(ptr))
	if err != nil {
		return 0, err
	}
	result := NSTime(int64(results[0]))
	if result == NSTimeError || result == NSTimeUnset {
		return result, fmt.Errorf("mseedio: parse time %q", value)
	}
	return result, nil
}

// FormatTime formats a libmseed nanosecond time without exposing WASM memory.
func (m *MiniSEED) FormatTime(value NSTime, format TimeFormat, subseconds SubsecondFormat) (string, error) {
	if err := m.ready(); err != nil {
		return "", err
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	const bufferSize = uint32(128)
	buffer, err := m.allocLocked(bufferSize)
	if err != nil {
		return "", err
	}
	defer m.freeLocked(buffer)

	results, err := m.callLocked(
		FnMSNSTimeToTimeStringN,
		uint64(value), uint64(buffer), uint64(bufferSize), uint64(format), uint64(subseconds),
	)
	if err != nil {
		return "", err
	}
	if uint32(results[0]) == 0 {
		return "", fmt.Errorf("mseedio: format time: libmseed returned NULL")
	}
	return m.readStringLocked(buffer)
}

func (m *MiniSEED) ready() error {
	if m == nil || !m.initialized || m.module == nil || m.memory == nil || m.Raw == nil {
		return fmt.Errorf("mseedio: runtime is not initialized")
	}
	return nil
}
