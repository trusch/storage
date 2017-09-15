package storage

import "io"

// WriteCloser wraps a writer and a closer
// on close first the writer is closed (if its an io.Closer)
// than the closer is closed
type WriteCloser struct {
	Closer io.Closer
	Writer io.Writer
}

func (s *WriteCloser) Write(p []byte) (int, error) {
	return s.Writer.(io.Writer).Write(p)
}

// Close closes the writer and the closer
func (s *WriteCloser) Close() error {
	if closer, ok := s.Writer.(io.Closer); ok {
		if err := closer.Close(); err != nil {
			return err
		}
	}
	return s.Closer.Close()
}
