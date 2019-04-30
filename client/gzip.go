package client

import (
	"compress/gzip"
	"io"
	"sync"
)

// gzipReaderPool contains all of our gzip writers.
// We use a pool so that every request can re-use writers.
var gzipReaderPool sync.Pool

// acquireGZipReader will return a clean gzip reader from the pool.
func acquireGZipReader(request io.Reader) (*gzip.Reader, error) {
	obj := gzipReaderPool.Get()

	if obj == nil {
		reader, err := gzip.NewReader(request)
		return reader, err
	}

	reader := obj.(*gzip.Reader)
	err := reader.Reset(request)
	return reader, err
}
