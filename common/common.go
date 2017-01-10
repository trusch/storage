package common

import "fmt"

// DocInfo describes a document
// It is used for directory listing
type DocInfo struct {
	Key   string
	Value []byte
}

// ListOpts are options given to the List command of storage implementations
// If Prefix != "" only docs with a key starting with the prefix are returned
// If Prefix == "" && Start != "" && End != "" only keys between Start and End are returned.
// Start and End are inclusive
type ListOpts struct {
	Prefix string
	Start  string
	End    string
}

// StorageError is the type of all possible storage errors
type StorageError struct {
	Msg    string
	Errors []error
}

// Error returns the string representation of this error
func (err *StorageError) Error() string {
	return fmt.Sprintf("%v: %v", err.Msg, err.Errors)
}

// StorageErrorType specifies the type of a storage related error
type StorageErrorType int

const (
	// BucketNotFound is thrown if the requested bucket is not yet created
	BucketNotFound StorageErrorType = iota
	// ReadFailed is thrown if the requested Key is not in the specified bucket
	ReadFailed
	// WriteFailed is thrown if the underlying db engine fail to write for some reason
	WriteFailed
	// CloseFailed is thrown if closing of the underlying db is not possible
	CloseFailed
	// InitFailed is thrown if opening of the underlying db is not possible
	InitFailed
)

// Error returns a StorageError with the specified type and info
func Error(typ StorageErrorType, errors ...error) error {
	switch typ {
	case BucketNotFound:
		return &StorageError{"bucket not found", errors}
	case ReadFailed:
		return &StorageError{"key not found", errors}
	case WriteFailed:
		return &StorageError{"write failed", errors}
	case CloseFailed:
		return &StorageError{"close failed", errors}
	case InitFailed:
		return &StorageError{"init failed", errors}
	}
	return &StorageError{"unknown storage error type", errors}
}
