package mseedio

import "fmt"

const (
	NoError        int32 = 0
	EndOfFile      int32 = 1
	GenericError   int32 = -1
	NotMiniSEED    int32 = -2
	WrongLength    int32 = -3
	OutOfRange     int32 = -4
	UnknownFormat  int32 = -5
	BadCompression int32 = -6
	InvalidCRC     int32 = -7
)

// LibraryError is an error code returned by libmseed.
type LibraryError struct {
	Operation string
	Code      int32
	Message   string
}

// NeedMoreDataError reports a valid record header whose payload is incomplete.
type NeedMoreDataError struct {
	Needed int32
}

func (e *NeedMoreDataError) Error() string {
	return fmt.Sprintf("mseedio: incomplete record: need %d more bytes", e.Needed)
}

func (e *LibraryError) Error() string {
	if e.Message == "" {
		return fmt.Sprintf("mseedio: %s: libmseed error %d", e.Operation, e.Code)
	}
	return fmt.Sprintf("mseedio: %s: %s (%d)", e.Operation, e.Message, e.Code)
}

func (m *MiniSEED) libraryErrorLocked(operation string, code int32) error {
	results, err := m.callLocked(FnMSErrorString, uint64(uint32(code)))
	if err != nil || len(results) == 0 || uint32(results[0]) == 0 {
		return &LibraryError{Operation: operation, Code: code}
	}
	message, readErr := m.readStringLocked(uint32(results[0]))
	if readErr != nil {
		return &LibraryError{Operation: operation, Code: code}
	}
	return &LibraryError{Operation: operation, Code: code, Message: message}
}
