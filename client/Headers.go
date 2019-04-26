package client

// Headers is a synonym for map[string]string.
type Headers map[string]string

// Common headers
var (
	headerEndSequence     = []byte{'\r', '\n', '\r', '\n'}
	contentLengthHeader   = []byte("Content-Length")
	contentEncodingHeader = []byte("Content-Encoding")
	gzipAccept            = []byte("gzip")
)
