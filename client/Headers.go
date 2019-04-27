package client

// Headers is a synonym for map[string]string.
type Headers map[string]string

// Common headers
var (
	newlineSequence        = []byte{'\r', '\n'}
	doubleNewlineSequence  = []byte{'\r', '\n', '\r', '\n'}
	contentLengthHeader    = []byte("Content-Length")
	contentEncodingHeader  = []byte("Content-Encoding")
	transferEncodingHeader = []byte("Transfer-Encoding")
	gzipAccept             = []byte("gzip")
	chunkedEncoding        = []byte("chunked")
)
