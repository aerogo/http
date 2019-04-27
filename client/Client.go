package client

import (
	"fmt"
	"log"
	"net"
	"net/url"

	jsoniter "github.com/json-iterator/go"
)

// Client represents a single
type Client struct {
	request  request
	response Response
}

// strings.Replace(path, " ", "%20", -1)

// Get builds a GET request.
func Get(path string) *Client {
	parsedURL, _ := url.Parse(path)

	http := &Client{
		request: request{
			method: "GET",
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

// Post builds a POST request.
func Post(path string) *Client {
	parsedURL, _ := url.Parse(path)

	http := &Client{
		request: request{
			method: "POST",
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
	data, err := jsoniter.Marshal(obj)

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

	for _, ip := range ips {
		err = http.exec(ip)

		// If it worked with one IP, we can stop here.
		// No need to test the other IPs.
		if err == nil {
			return nil
		}
	}

	return err
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
