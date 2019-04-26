package client

import (
	"bytes"
)

// var gzipType = []byte("gzip")

// Response represents the HTTP response used in the given client.
type Response struct {
	header []byte
	body   []byte
}

// StatusCode returns the status code of the response.
func (response Response) StatusCode() int {
	return 0
}

// String returns the response body as a string.
func (response Response) String() string {
	return string(response.body)
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

// // String returns the response body as a string.
// func (response Response) String() string {
// 	return string(response.Bytes())
// }

// // RawString returns the raw response body as a string.
// func (response Response) RawString() string {
// 	return string(response.RawBytes())
// }

// // Bytes returns the response body as a byte slice and unzips gzipped content when necessary.
// func (response Response) Bytes() []byte {
// 	encoding := response.inner.Header.Peek("Content-Encoding")

// 	if bytes.Equal(encoding, gzipType) {
// 		unzipped, err := response.inner.BodyGunzip()

// 		if err != nil {
// 			return response.inner.Body()
// 		}

// 		return unzipped
// 	}

// 	return response.inner.Body()
// }

// // RawBytes returns the raw response body as a byte slice.
// func (response Response) RawBytes() []byte {
// 	return response.inner.Body()
// }

// // Unmarshal tries to JSON decode the response and save it in the object.
// func (response Response) Unmarshal(obj interface{}) error {
// 	return jsoniter.Unmarshal(response.Bytes(), obj)
// }
