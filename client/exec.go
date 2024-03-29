package client

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net"
	"unicode"

	"github.com/aerogo/http/ciphers"
	"github.com/akyoto/stringutils/convert"
)

var clientSessionCache = tls.NewLRUClientSessionCache(0)

// exec executes the request and returns the response.
func (http *Client) exec(connection net.Conn) error {
	err := connection.(*net.TCPConn).SetNoDelay(true)

	if err != nil {
		return err
	}

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
	var (
		response          bytes.Buffer
		decodedChunks     bytes.Buffer
		tmp               = make([]byte, 16384)
		headerEndPosition = -1
		bodyStartPosition = -1
		lastChunkPosition = -1
		contentLength     = -1
		isChunked         = false
	)

	for {
		n, err := connection.Read(tmp)
		response.Write(tmp[:n])

		// Find status
		if http.response.statusCode == 0 {
			statusPos := bytes.IndexByte(tmp, ' ')
			statusSlice := tmp[statusPos+1 : statusPos+4]
			http.response.statusCode = int(convert.DecToInt(statusSlice))
		}

		// Find end of headers
		if headerEndPosition == -1 {
			headerEndPosition = bytes.Index(response.Bytes(), doubleNewlineSequence)

			if headerEndPosition != -1 {
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
					contentLength = int(convert.DecToInt(lengthSlice))
					response.Grow(contentLength)
				}
			}
		}

		if isChunked {
			chunkBytes := response.Bytes()[lastChunkPosition:]
			chunkBytesRead, finished := decodeChunks(chunkBytes, &decodedChunks)

			if finished {
				http.response.body = decodedChunks.Bytes()
				return nil
			}

			lastChunkPosition += chunkBytesRead
		} else if headerEndPosition != -1 && response.Len()-bodyStartPosition >= contentLength {
			http.response.body = response.Bytes()[bodyStartPosition:]
			return nil
		}

		if err != nil {
			http.response.body = response.Bytes()[bodyStartPosition:]
			return err
		}
	}
}
