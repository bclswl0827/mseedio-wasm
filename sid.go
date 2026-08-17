package mseedio

import (
	"fmt"
	"strings"
)

// NSLC contains the FDSN network, station, location and channel codes.
type NSLC struct {
	Network  string
	Station  string
	Location string
	Channel  string
}

// SIDToNSLC splits an FDSN source identifier into its component codes.
func (m *MiniSEED) SIDToNSLC(sid string) (NSLC, error) {
	if err := m.ready(); err != nil {
		return NSLC{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	sidPtr, err := m.allocBytesLocked(append([]byte(sid), 0))
	if err != nil {
		return NSLC{}, err
	}
	defer m.freeLocked(sidPtr)

	const componentSize = uint32(64)
	pointers := make([]uint32, 4)
	for index := range pointers {
		pointers[index], err = m.allocLocked(componentSize)
		if err != nil {
			for _, ptr := range pointers[:index] {
				_ = m.freeLocked(ptr)
			}
			return NSLC{}, err
		}
		defer m.freeLocked(pointers[index])
	}

	results, err := m.callLocked(
		FnMSSidToNSLCN,
		uint64(sidPtr),
		uint64(pointers[0]), uint64(componentSize),
		uint64(pointers[1]), uint64(componentSize),
		uint64(pointers[2]), uint64(componentSize),
		uint64(pointers[3]), uint64(componentSize),
	)
	if err != nil {
		return NSLC{}, err
	}
	if code := int32(results[0]); code != NoError {
		return NSLC{}, m.libraryErrorLocked("split source identifier", code)
	}

	values := make([]string, 4)
	for index, ptr := range pointers {
		values[index], err = m.readStringLocked(ptr)
		if err != nil {
			return NSLC{}, err
		}
	}
	return NSLC{Network: values[0], Station: values[1], Location: values[2], Channel: values[3]}, nil
}

// NSLCToSID creates a source identifier from FDSN component codes.
func (m *MiniSEED) NSLCToSID(nslc NSLC) (string, error) {
	if err := m.ready(); err != nil {
		return "", err
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	const sidSize = uint32(64)
	sidPtr, err := m.allocLocked(sidSize)
	if err != nil {
		return "", err
	}
	defer m.freeLocked(sidPtr)

	components := []string{nslc.Network, nslc.Station, nslc.Location, nslc.Channel}
	pointers := make([]uint32, len(components))
	for index, component := range components {
		pointers[index], err = m.allocBytesLocked(append([]byte(component), 0))
		if err != nil {
			for _, ptr := range pointers[:index] {
				_ = m.freeLocked(ptr)
			}
			return "", err
		}
		defer m.freeLocked(pointers[index])
	}

	results, err := m.callLocked(
		FnMSNSLCToSid,
		uint64(sidPtr), uint64(sidSize), 0,
		uint64(pointers[0]), uint64(pointers[1]), uint64(pointers[2]), uint64(pointers[3]),
	)
	if err != nil {
		return "", err
	}
	if code := int32(results[0]); code < NoError {
		return "", m.libraryErrorLocked("create source identifier", code)
	}
	return m.readStringLocked(sidPtr)
}

// NormalizeSID accepts either FDSN:NET_STA_LOC_B_S_SS or NET.STA.LOC.BSS.
func (m *MiniSEED) NormalizeSID(value string) (string, error) {
	if strings.HasPrefix(value, "FDSN:") {
		if _, err := m.SIDToNSLC(value); err != nil {
			return "", err
		}
		return value, nil
	}
	parts := strings.Split(value, ".")
	if len(parts) != 4 {
		return "", fmt.Errorf("mseedio: NSLC must be NET.STA.LOC.CHAN")
	}
	return m.NSLCToSID(NSLC{Network: parts[0], Station: parts[1], Location: parts[2], Channel: parts[3]})
}
