package extractor

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

var webpHeader = []byte("WEBP")

// Skip RIFF and File Size header
var skipBytesCount = 4 + 4

type WEBP struct{}

func (e WEBP) BufSize() int {
	return skipBytesCount + len(webpHeader)
}

func (e WEBP) MatchFormat(buf []byte) (string, bool) {
	return "webp", bytes.Equal(buf[skipBytesCount:e.BufSize()], webpHeader)
}

func (e WEBP) ExtractSize(reader io.ReadSeeker) (width, height int, err error) {
	var buffer [4]byte

	var skipBytes int64 = int64(skipBytesCount + len(webpHeader))
	if _, err = reader.Seek(skipBytes, io.SeekStart); err != nil {
		err = fmt.Errorf("failed to seek to correct position: %w", err)
		return
	}

	// Get the VP8 tag
	if _, err = reader.Read(buffer[:]); err != nil {
		err = fmt.Errorf("failed to read buffer: %w", err)
		return
	}

	switch buffer[3] {
	case ' ':
		return e.webpVp8Size(reader)
	case 'L':
		return e.webpVp8lSize(reader)
	case 'X':
		return e.webpVp8xSize(reader)
	default:
		err = errors.New("unknown VP8 tag")
		return
	}
}

func (e WEBP) webpVp8Size(reader io.ReadSeeker) (width, height int, err error) {
	if _, err = reader.Seek(26, io.SeekStart); err != nil {
		err = fmt.Errorf("failed to seek: %w", err)
		return
	}

	widthU16, widthErr := readU16(reader, LittleEndian)
	heightU16, heightErr := readU16(reader, LittleEndian)
	return int(widthU16), int(heightU16), errors.Join(widthErr, heightErr)
}

func (e WEBP) webpVp8lSize(reader io.ReadSeeker) (width, height int, err error) {
	if _, err = reader.Seek(21, io.SeekStart); err != nil {
		err = fmt.Errorf("failed to seek: %w", err)
		return
	}

	dims, err := readU32(reader, LittleEndian)
	if err != nil {
		err = fmt.Errorf("failed to read dimensions: %w", err)
		return
	}

	// Extract the width and height from the 32-bit integer (packed format).
	width = int(dims & 0x3FFF)
	height = int((dims >> 14) & 0x3FFF)
	return
}

func (e WEBP) webpVp8xSize(reader io.ReadSeeker) (width, height int, err error) {
	if _, err = reader.Seek(24, io.SeekStart); err != nil {
		err = fmt.Errorf("failed to seek: %w", err)
		return
	}

	widthU24, widthErr := readU24(reader, LittleEndian)
	heightU24, heightErr := readU24(reader, LittleEndian)
	return int(widthU24), int(heightU24), errors.Join(widthErr, heightErr)
}
