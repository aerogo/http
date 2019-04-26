package client

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/url"
	"strconv"
	"strings"

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
				"Accept-Encoding": "gzip",
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
				"Accept-Encoding": "gzip",
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
func (http *Client) Body(raw string) *Client {
	http.request.body = raw
	return http
}

// BodyJSON sets the request body by converting the object to JSON.
func (http *Client) BodyJSON(obj interface{}) *Client {
	data, err := jsoniter.MarshalToString(obj)

	if err != nil {
		log.Printf("Error converting request body to JSON: %v", err)
		return http
	}

	http.request.body = data
	return http
}

// BodyBytes sets the request body as a byte slice.
func (http *Client) BodyBytes(raw []byte) *Client {
	http.request.body = string(raw)
	return http
}

// Response returns the response object.
func (http *Client) Response() *Response {
	return &http.response
}

// Do executes the request and returns the response.
func (http *Client) Do() error {
	var connection net.Conn
	var err error

	hostName := http.request.url.Hostname()
	port, _ := strconv.Atoi(http.request.url.Port())
	path := http.request.url.Path

	if port == 0 {
		if http.request.url.Scheme == "https" {
			port = 443
		} else {
			port = 80
		}
	}

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	if http.request.url.Scheme == "https" {
		// TLS
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
		}

		connection, err = tls.Dial("tcp", fmt.Sprintf("%s:%d", hostName, port), tlsConfig)

		if err != nil {
			return err
		}
	} else {
		ips, err := net.LookupIP(hostName)

		if err != nil {
			return err
		}

		if len(ips) == 0 {
			return fmt.Errorf("Could not resolve host: %s", hostName)
		}

		remoteAddress := net.TCPAddr{
			IP:   ips[0],
			Port: port,
		}

		connection, err = net.DialTCP("tcp", nil, &remoteAddress)

		if err != nil {
			return err
		}

		connection.(*net.TCPConn).SetNoDelay(true)
	}

	defer connection.Close()

	// tlsClient := tls.Client(connection, tlsConfig)
	// err = tlsClient.Handshake()

	// if err != nil {
	// 	return err
	// }

	// Create request headers
	var requestHeaders bytes.Buffer

	requestHeaders.WriteString("GET ")
	requestHeaders.WriteString(path)
	requestHeaders.WriteString(" HTTP/1.1\r\nHost: ")
	requestHeaders.WriteString(hostName)
	requestHeaders.WriteString("\r\n")

	for key, value := range http.request.headers {
		requestHeaders.WriteString(key)
		requestHeaders.WriteString(": ")
		requestHeaders.WriteString(value)
		requestHeaders.WriteString("\r\n")
	}

	requestHeaders.WriteString("\r\n")

	connection.Write(requestHeaders.Bytes())

	var header bytes.Buffer
	var body bytes.Buffer
	current := &header
	tmp := make([]byte, 16384)
	contentLength := 0

	for {
		n, err := connection.Read(tmp)
		headerEnd := bytes.Index(tmp, headerEndSequence)

		if headerEnd != -1 {
			header.Write(tmp[:headerEnd])
			body.Write(tmp[headerEnd+4 : n])
			current = &body

			// Find content length
			http.response.header = header.Bytes()
			lengthSlice := http.response.Header(contentLengthHeader)
			contentLength = asciiToInt(lengthSlice)

			// Find status
			statusPos := bytes.IndexByte(http.response.header, ' ')
			statusSlice := http.response.header[statusPos+1 : statusPos+4]
			http.response.statusCode = asciiToInt(statusSlice)

			// Reserve space for the content length
			body.Grow(contentLength)
		} else {
			current.Write(tmp[:n])
		}

		if err != nil || body.Len() >= contentLength {
			http.response.body = body.Bytes()
			return err
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
