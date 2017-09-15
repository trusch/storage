package storage

import "io"

// ReadCloser wraps a closer and a reader.
// When closing first reader is closed if its a io.Closer
// Then the closer is closed
type ReadCloser struct {
	Closer io.Closer
	Reader io.Reader
}

func (s *ReadCloser) Read(p []byte) (int, error) {
	return s.Reader.Read(p)
}

// Close closes the wrapped reader and closer.
func (s *ReadCloser) Close() error {
	if closer, ok := s.Reader.(io.Closer); ok {
		if err := closer.Close(); err != nil {
			return err
		}
	}
	return s.Closer.Close()
}
