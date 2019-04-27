package client

import (
	"bytes"
	"io/ioutil"

	jsoniter "github.com/json-iterator/go"
)

// Response represents the HTTP response used in the given client.
type Response struct {
	header     []byte
	body       []byte
	statusCode int
}

// StatusCode returns the status code of the response.
func (response Response) StatusCode() int {
	return response.statusCode
}

// Ok returns true if the HTTP status code is not lower than 200 and not 400 or higher.
func (response Response) Ok() bool {
	return response.statusCode >= 200 && response.statusCode < 400
}

// Header returns the value for the given header.
func (response Response) Header(name []byte) []byte {
	start := bytes.Index(response.header, append(newlineSequence, name...))

	if start == -1 {
		return nil
	}

	start += len(newlineSequence) + len(name) + len(": ")
	remaining := response.header[start:]
	end := bytes.IndexByte(remaining, '\r')

	if end == -1 {
		return remaining
	}

	return remaining[:end]
}

// HeaderString returns the string value for the given header.
func (response Response) HeaderString(name string) string {
	return string(response.Header([]byte(name)))
}

// RawHeaders returns headers as a byte slice.
func (response Response) RawHeaders() []byte {
	return response.header
}

// RawHeadersString returns headers as a string.
func (response Response) RawHeadersString() string {
	return string(response.header)
}

// Bytes returns the response body as a byte slice and unzips gzipped content when necessary.
func (response Response) Bytes() []byte {
	encoding := response.Header(contentEncodingHeader)

	if !bytes.Equal(encoding, gzipAccept) {
		return response.body
	}

	bodyReader := bytes.NewReader(response.body)
	gzipReader := acquireGZipReader(bodyReader)
	unzipped, err := ioutil.ReadAll(gzipReader)

	if err != nil {
		return response.body
	}

	return unzipped
}

// String returns the response body as a string.
func (response Response) String() string {
	return string(response.Bytes())
}

// RawBytes returns the raw response body as a byte slice.
func (response Response) RawBytes() []byte {
	return response.body
}

// RawString returns the raw response body as a string.
func (response Response) RawString() string {
	return string(response.body)
}

// Unmarshal tries to JSON decode the response and save it in the object.
func (response Response) Unmarshal(obj interface{}) error {
	return jsoniter.Unmarshal(response.Bytes(), obj)
}
