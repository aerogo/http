package client

import (
	"strings"

	"github.com/valyala/fasthttp"
)

// Headers is a synonym for map[string]string.
type Headers map[string]string

// Client ...
type Client struct {
	client   fasthttp.Client
	request  *fasthttp.Request
	response *fasthttp.Response
}

// Get builds a GET request.
func Get(path string) *Client {
	http := new(Client)
	http.request = fasthttp.AcquireRequest()
	http.response = fasthttp.AcquireResponse()

	http.request.SetRequestURI(strings.Replace(path, " ", "%20", -1))
	http.request.Header.Set("Accept-Encoding", "gzip")
	return http
}

// Post builds a POST request.
func Post(path string) *Client {
	http := new(Client)
	http.request = fasthttp.AcquireRequest()
	http.response = fasthttp.AcquireResponse()

	http.request.SetRequestURI(strings.Replace(path, " ", "%20", -1))
	http.request.Header.SetMethod("POST")
	return http
}

// Header sets one HTTP header for the request.
func (http *Client) Header(key string, value string) *Client {
	http.request.Header.Set(key, value)
	return http
}

// Headers sets the HTTP headers for the request.
func (http *Client) Headers(headers Headers) *Client {
	for key, value := range headers {
		http.request.Header.Set(key, value)
	}
	return http
}

// Body sets the request body.
func (http *Client) Body(raw string) *Client {
	http.request.SetBodyString(raw)
	return http
}

// BodyBytes sets the request body as a byte slice.
func (http *Client) BodyBytes(raw []byte) *Client {
	http.request.SetBody(raw)
	return http
}

// Response returns the response object.
func (http *Client) Response() Response {
	return Response{
		inner: http.response,
	}
}

// End executes the request and returns the response.
func (http *Client) End() (Response, error) {
	err := http.client.Do(http.request, http.response)
	return http.Response(), err
}

// EndStruct executes the request, unmarshals the response body into a struct and returns the response.
func (http *Client) EndStruct(obj interface{}) (Response, error) {
	err := http.client.Do(http.request, http.response)
	response := http.Response()

	if err != nil {
		return response, err
	}

	err = response.Unmarshal(obj)
	return response, err
}
