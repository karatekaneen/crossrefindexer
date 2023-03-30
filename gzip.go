package crossrefindexer

import (
	"compress/gzip"
	"fmt"
	"io"
)

// GzipReader creates a new Reader of a gziped file.
//
// It is the caller's responsibility to call Close on the Reader when done.
func GzipReader(r io.Reader) (*gzip.Reader, error) {
	gzipReader, err := gzip.NewReader(r)
	if err != nil {
		return gzipReader, fmt.Errorf("create gzip reader: %w", err)
	}
	return gzipReader, nil
}
