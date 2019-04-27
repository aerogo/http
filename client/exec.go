package client

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/aerogo/http/ciphers"
	"net"
	"strconv"
	"unicode"

	"github.com/aerogo/http/convert"
)

var clientSessionCache = tls.NewLRUClientSessionCache(0)

// exec executes the request and returns the respense for the given IP.
func (http *Client) exec(ip net.IP) error {
	var connection net.Conn
	var err error

	port, _ := strconv.Atoi(http.request.url.Port())

	if port == 0 {
		if http.request.url.Scheme == "https" {
			port = 443
		} else {
			port = 80
		}
	}

	remoteAddress := net.TCPAddr{
		IP:   ip,
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
			ServerName:         http.request.url.Hostname(),
			MinVersion:         tls.VersionTLS12,
			MaxVersion:         tls.VersionTLS13,
			CipherSuites:       ciphers.List,
			ClientSessionCache: clientSessionCache,
			InsecureSkipVerify: true,
		}

		connection = tls.Client(connection, tlsConfig)
	}

	// Make sure we close the connection
	defer connection.Close()

	// Create request headers
	var requestContents bytes.Buffer

	requestContents.WriteString(http.request.method)
	requestContents.WriteByte(' ')
	requestContents.WriteString(http.request.url.RequestURI())
	requestContents.WriteString(" HTTP/1.1\r\n")

	for key, value := range http.request.headers {
		requestContents.WriteString(key)
		requestContents.WriteString(": ")
		requestContents.WriteString(value)
		requestContents.WriteString("\r\n")
	}

	if len(http.request.body) > 0 {
		fmt.Fprintf(&requestContents, "Content-Length: %d\r\n", len(http.request.body))
	}

	requestContents.WriteString("\r\n")
	requestContents.Write(http.request.body)

	// Send request headers
	_, err = connection.Write(requestContents.Bytes())

	if err != nil {
		return err
	}

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
			// Find status
			if http.response.statusCode == 0 {
				statusPos := bytes.IndexByte(tmp, ' ')
				statusSlice := tmp[statusPos+1 : statusPos+4]
				http.response.statusCode = convert.ASCIIDecToInt(statusSlice)
			}

			doubleNewlinePos := bytes.Index(tmp, doubleNewlineSequence)

			if doubleNewlinePos != -1 {
				headerEndPosition = response.Len() - n + doubleNewlinePos
				bodyStartPosition = headerEndPosition + len(doubleNewlineSequence)
				lastChunkPosition = bodyStartPosition
				http.response.header = response.Bytes()[:headerEndPosition]

				// Normalize headers
				makeUpper := true
				isValue := false

				for index, char := range http.response.header {
					if char == '\n' {
						isValue = false
						makeUpper = true
						continue
					}

					if char == '-' {
						makeUpper = true
						continue
					}

					if char == ':' {
						isValue = true
						continue
					}

					if makeUpper && !isValue {
						http.response.header[index] = byte(unicode.ToUpper(rune(char)))
						makeUpper = false
					}
				}

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
