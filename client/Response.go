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

// Header returns the value for the given header.
func (response Response) Header(name []byte) []byte {
	start := bytes.Index(response.header, name)

	if start == -1 {
		return nil
	}

	start += len(name) + 2
	remaining := response.header[start:]
	end := bytes.IndexByte(remaining, '\r')

	if end == -1 {
		return nil
	}

	return remaining[:end]
}

// HeaderString returns the string value for the given header.
func (response Response) HeaderString(name string) string {
	return string(response.Header([]byte(name)))
}

// Bytes returns the response body as a byte slice and unzips gzipped content when necessary.
func (response Response) Bytes() []byte {
	encoding := response.Header(contentEncodingHeader)

	if bytes.Equal(encoding, gzipAccept) {
		bodyReader := bytes.NewReader(response.body)
		gzipReader := acquireGZipReader(bodyReader)
		unzipped, err := ioutil.ReadAll(gzipReader)

		if err != nil {
			return response.body
		}

		return unzipped
	}

	return response.body
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
