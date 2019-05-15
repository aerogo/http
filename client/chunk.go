package client

import (
	"bytes"

	"github.com/akyoto/stringutils/convert"
)

func decodeChunks(raw []byte, target *bytes.Buffer) (int, bool) {
	bytesRead := 0

	for {
		newlinePosition := bytes.Index(raw, newlineSequence)

		if newlinePosition == -1 {
			return bytesRead, false
		}

		expectedChunkBodyLength := convert.HexToInt(raw[:newlinePosition])

		// A chunk with 0 length indicates that we're finished.
		if expectedChunkBodyLength == 0 {
			return bytesRead, true
		}

		chunkBodyStart := newlinePosition + len(newlineSequence)
		chunkBodyEnd := chunkBodyStart + expectedChunkBodyLength
		willReadBytes := chunkBodyEnd + len(newlineSequence)

		// If the chunk didn't fully arrive yet...
		if len(raw) < willReadBytes {
			// ...do nothing. Parse it on the next iteration.
			return bytesRead, false
		}

		currentChunk := raw[chunkBodyStart:chunkBodyEnd]
		target.Write(currentChunk)
		bytesRead += willReadBytes
		raw = raw[willReadBytes:]
	}
}
