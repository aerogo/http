package client

import "net/url"

// request represents the HTTP request used in the given client.
type request struct {
	method  string
	url     *url.URL
	body    []byte
	headers Headers
}
