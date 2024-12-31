package extractor

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/pillowskiy/imagesize/imagebytes"
)

var (
	pngHeader  = []byte("\x89\x50\x4E\x47")
	ihdrHeader = []byte("IHDR")
)

// PNG defines an extractor for PNG image format.
//
// The PNG file format starts with a specific header and chunk structure:
//
// 1. The first 8 bytes contain the following:
//   - 0x89 0x50 0x4E 0x47: ASCII characters "PNG" with a special signature to identify the file as a PNG file.
//   - 0x0D 0x0A 0x1A 0x0A: A sequence of special byte values indicating the start of the PNG file.
//
// 2. The next 4 bytes represent the length of the IHDR chunk (13 bytes).
// 3. The next 4 bytes contain the ASCII characters "IHDR", identifying the chunk as the image header.
// 4. The next 4 bytes represent the width of the image (1 pixel in this case), encoded as a 32-bit unsigned integer in big-endian format.
// 5. The next 4 bytes represent the height of the image (1 pixel in this case), also encoded as a 32-bit unsigned integer in big-endian format.
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

	widthU32, widthErr := imagebytes.ReadU32(reader, imagebytes.BigEndian)
	heightU32, heightErr := imagebytes.ReadU32(reader, imagebytes.BigEndian)
	return int(widthU32), int(heightU32), errors.Join(widthErr, heightErr)
}
