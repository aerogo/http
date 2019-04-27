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

	"github.com/aerogo/http/convert"
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

	if http.request.url.Scheme == "https" {
		// TLS
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
		}

		connection = tls.Client(connection, tlsConfig)
	}

	// Make sure we close the connection
	defer connection.Close()

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

	// Send request
	connection.Write(requestHeaders.Bytes())

	// Receive response
	var response bytes.Buffer
	tmp := make([]byte, 16384)
	contentLength := 0
	isChunked := false
	headerEndPosition := -1
	bodyStartPosition := 0
	lastChunkPosition := -1

	// Create another buffer just for chunked responses
	var decodedChunks bytes.Buffer

	for {
		n, err := connection.Read(tmp)
		response.Write(tmp[:n])

		// Find headers
		if headerEndPosition == -1 {
			doubleNewlinePos := bytes.Index(tmp, doubleNewlineSequence)

			if doubleNewlinePos != -1 {
				headerEndPosition = response.Len() - n + doubleNewlinePos
				bodyStartPosition = headerEndPosition + len(doubleNewlineSequence)
				lastChunkPosition = bodyStartPosition
				http.response.header = response.Bytes()[:headerEndPosition]

				// Find status
				statusPos := bytes.IndexByte(http.response.header, ' ')
				statusSlice := http.response.header[statusPos+1 : statusPos+4]
				http.response.statusCode = convert.ASCIIDecToInt(statusSlice)

				// Find content length
				transferSlice := http.response.Header(transferEncodingHeader)

				if bytes.Equal(transferSlice, chunkedEncoding) {
					isChunked = true
				} else {
					lengthSlice := http.response.Header(contentLengthHeader)
					contentLength = convert.ASCIIDecToInt(lengthSlice)
					response.Grow(contentLength)
				}
			}
		}

		// Read chunks
		if isChunked {
			chunkBytes := response.Bytes()[lastChunkPosition:]
			chunkBytesRead, finished := decodeChunks(chunkBytes, &decodedChunks)

			// If we read the last chunk, we're finished here
			if finished {
				http.response.body = decodedChunks.Bytes()
				return nil
			}

			lastChunkPosition += chunkBytesRead
		}

		// End response on error
		if err != nil {
			http.response.body = response.Bytes()[bodyStartPosition:]
			return err
		}

		// End response if content length has been reached
		if !isChunked && response.Len()-bodyStartPosition >= contentLength {
			http.response.body = response.Bytes()[bodyStartPosition:]
			return nil
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
