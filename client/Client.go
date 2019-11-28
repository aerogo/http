package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/url"
	"strconv"
)

// Client represents a single
type Client struct {
	request  request
	response Response
}

// Get builds a GET request.
func Get(path string) *Client {
	return WithMethod(path, "GET")
}

// Post builds a POST request.
func Post(path string) *Client {
	return WithMethod(path, "POST")
}

// Put builds a PUT request.
func Put(path string) *Client {
	return WithMethod(path, "PUT")
}

// Delete builds a DELETE request.
func Delete(path string) *Client {
	return WithMethod(path, "DELETE")
}

// WithMethod builds a request with a custom HTTP verb.
func WithMethod(path string, method string) *Client {
	parsedURL, err := url.Parse(path)

	if err != nil {
		panic("Invalid URL: " + path)
	}

	http := &Client{
		request: request{
			method: method,
			url:    parsedURL,
			headers: Headers{
				"Host":            parsedURL.Hostname(),
				"Accept-Encoding": "gzip",
				"Accept":          "*/*",
			},
		},
	}

	return http
}

// Header sets one HTTP header for the request.
func (http *Client) Header(key string, value string) *Client {
	http.request.headers[key] = value
	return http
}

// Headers sets the HTTP headers for the request.
func (http *Client) Headers(headers Headers) *Client {
	for key, value := range headers {
		http.request.headers[key] = value
	}

	return http
}

// Body sets the request body.
func (http *Client) Body(raw []byte) *Client {
	http.request.body = raw
	return http
}

// BodyString sets the request body.
func (http *Client) BodyString(raw string) *Client {
	http.request.body = []byte(raw)
	return http
}

// BodyJSON sets the request body by converting the object to JSON.
func (http *Client) BodyJSON(obj interface{}) *Client {
	data, err := json.Marshal(obj)

	if err != nil {
		log.Printf("Error converting request body to JSON: %v", err)
		return http
	}

	http.request.body = data
	return http
}

// Response returns the response object.
func (http *Client) Response() *Response {
	return &http.response
}

// Do executes the request and returns the response.
func (http *Client) Do() error {
	ips, err := net.LookupIP(http.request.url.Hostname())

	if err != nil {
		return err
	}

	if len(ips) == 0 {
		return fmt.Errorf("Could not resolve host: %s", http.request.url.Hostname())
	}

	port, _ := strconv.Atoi(http.request.url.Port())

	if port == 0 {
		if http.request.url.Scheme == "https" {
			port = 443
		} else {
			port = 80
		}
	}

	connections := make(chan net.Conn, len(ips))
	errors := make(chan error, len(ips))

	for _, ip := range ips {
		go func(ip net.IP) {
			remoteAddress := net.TCPAddr{
				IP:   ip,
				Port: port,
			}

			connection, err := net.DialTCP("tcp", nil, &remoteAddress)

			if err != nil {
				errors <- err
				return
			}

			connections <- connection
		}(ip)
	}

	errorCount := 0

	for {
		select {
		case connection := <-connections:
			return http.exec(connection)

		case err := <-errors:
			errorCount++

			if errorCount == len(ips) {
				return err
			}
		}
	}
}

// End executes the request and returns the response.
func (http *Client) End() (*Response, error) {
	err := http.Do()
	return &http.response, err
}

// EndStruct executes the request, unmarshals the response body into a struct and returns the response.
func (http *Client) EndStruct(obj interface{}) (*Response, error) {
	err := http.Do()

	if err != nil {
		return &http.response, err
	}

	err = http.response.Unmarshal(obj)
	return &http.response, err
}
