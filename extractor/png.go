package extractor

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

var (
	pngHeader  = []byte("\x89\x50\x4E\x47")
	ihdrHeader = []byte("IHDR")
)

type PNG struct{}

func (e PNG) BufSize() int {
	return len(pngHeader)
}

func (e PNG) MatchFormat(buf []byte) (string, bool) {
	return "png", bytes.HasPrefix(buf, pngHeader)
}

func (e PNG) ExtractSize(reader io.ReadSeeker) (width, height int, err error) {
	if _, err = reader.Seek(12, io.SeekStart); err != nil {
		err = fmt.Errorf("failed to seek to the start of the file: %w", err)
		return
	}

	var buf [4]byte
	// Read potential IHDR header (12:16)
	_, err = io.ReadFull(reader, buf[:])
	if err != nil {
		err = fmt.Errorf("failed to read the first 4 bytes: %w", err)
		return
	}

	var skipHeaderBytesCount int64 = 8
	if bytes.Equal(buf[:], ihdrHeader) {
		skipHeaderBytesCount = 16
	}

	if _, err = reader.Seek(skipHeaderBytesCount, io.SeekStart); err != nil {
		err = fmt.Errorf("failed to seek to correct position: %w", err)
		return
	}

	widthU32, widthErr := readU32(reader, BigEndian)
	heightU32, heightErr := readU32(reader, BigEndian)
	return int(widthU32), int(heightU32), errors.Join(widthErr, heightErr)
}
