package client

import (
	"bytes"
	"encoding/json"

	"github.com/valyala/fasthttp"
)

var gzipType = []byte("gzip")

// Response represents the HTTP response used in the given client.
type Response struct {
	inner *fasthttp.Response
}

// StatusCode returns the status code of the response.
func (response Response) StatusCode() int {
	return response.inner.StatusCode()
}

// Body returns the response body as a string.
func (response Response) Body() string {
	return string(response.BodyBytes())
}

// RawBody returns the raw response body as a string.
func (response Response) RawBody() string {
	return string(response.RawBodyBytes())
}

// BodyBytes returns the response body as a byte slice and unzips gzipped content when necessary.
func (response Response) BodyBytes() []byte {
	encoding := response.inner.Header.Peek("Content-Encoding")

	if bytes.Equal(encoding, gzipType) {
		unzipped, err := response.inner.BodyGunzip()

		if err != nil {
			return response.inner.Body()
		}

		return unzipped
	}

	return response.inner.Body()
}

// RawBodyBytes returns the raw response body as a byte slice.
func (response Response) RawBodyBytes() []byte {
	return response.inner.Body()
}

// Unmarshal tries to JSON decode the response and save it in the object.
func (response Response) Unmarshal(obj interface{}) error {
	return json.Unmarshal(response.BodyBytes(), obj)
}
