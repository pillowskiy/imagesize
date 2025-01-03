package extractor

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/pillowskiy/imagesize/imagebytes"
)

var webpHeader = []byte("WEBP")

// Skip RIFF and File Size headers
var skipBytesCount = 4 + 4

// WEBP defines an extractor for WebP image format.
//
// The WebP file format starts with a RIFF-based structure:
// 1. The first 4 bytes contain the ASCII characters "RIFF", which identify the file as a RIFF container.
// 2. The next 4 bytes represent the file size (unsigned 32-bit integer, little-endian). This value is the total size of the file minus 8 bytes (header and size field).
// 3. The following 4 bytes contain the ASCII characters "WEBP", identifying the format as WebP.
// 4. The next 4 bytes specify the WebP encoding format ("VP8 ", "VP8L", "VP8X")
//
// Together, these 16 bytes form the mandatory WebP header.
type WEBP struct{}

func (e WEBP) BufSize() int {
	return skipBytesCount + len(webpHeader)
}

func (e WEBP) MatchFormat(buf []byte) (string, bool) {
	return "webp", bytes.Equal(buf[skipBytesCount:e.BufSize()], webpHeader)
}

func (e WEBP) ExtractSize(reader io.ReadSeeker) (width, height int, err error) {
	// Skip RIFF, FileSize and WEBP headers
	var skipBytes int64 = int64(skipBytesCount + len(webpHeader))
	if _, err = reader.Seek(skipBytes, io.SeekStart); err != nil {
		err = fmt.Errorf("failed to seek to correct position: %w", err)
		return
	}

	var buffer [4]byte
	// Read the VP8 Tag ("VP8 ", "VP8L", "VP8X")
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

// The WebP VP8 format stores the width and height as unsigned 16-bit integers at byte offsets 26 and 28
// in the image data (little-endian format). The function seeks to the appropriate position in the file,
// reads the width and height, and returns them.
func (e WEBP) webpVp8Size(reader io.ReadSeeker) (width, height int, err error) {
	if _, err = reader.Seek(26, io.SeekStart); err != nil {
		err = fmt.Errorf("failed to seek: %w", err)
		return
	}

	widthU16, widthErr := imagebytes.ReadU16(reader, imagebytes.LittleEndian)
	heightU16, heightErr := imagebytes.ReadU16(reader, imagebytes.LittleEndian)
	return int(widthU16), int(heightU16), errors.Join(widthErr, heightErr)
}

// The WebP VP8L format stores the width and height as packed data in a single 32-bit integer at byte offset 21.
// The lower 14 bits represent the width and the upper 14 bits represent the height (little-endian format).
// The function seeks to the appropriate position in the file, reads the 32-bit integer, and extracts the width and height.
// See https://developers.google.com/speed/webp/docs/webp_lossless_bitstream_specification
func (e WEBP) webpVp8lSize(reader io.ReadSeeker) (width, height int, err error) {
	if _, err = reader.Seek(21, io.SeekStart); err != nil {
		err = fmt.Errorf("failed to seek: %w", err)
		return
	}

	dims, err := imagebytes.ReadU32(reader, imagebytes.LittleEndian)
	if err != nil {
		err = fmt.Errorf("failed to read dimensions: %w", err)
		return
	}

	// Extract the width and height from the 32-bit integer (packed format).
	width = int(dims&0x3FFF) + 1
	height = int((dims>>14)&0x3FFF) + 1
	return
}

// The WebP VP8X format stores the width and height as unsigned 24-bit integers at byte offsets 24 and 27
// in the image data (little-endian format). The function seeks to the appropriate position in the file,
// reads the width and height, and returns them.
func (e WEBP) webpVp8xSize(reader io.ReadSeeker) (width, height int, err error) {
	if _, err = reader.Seek(24, io.SeekStart); err != nil {
		err = fmt.Errorf("failed to seek: %w", err)
		return
	}

	widthU24, widthErr := imagebytes.ReadU24(reader, imagebytes.LittleEndian)
	heightU24, heightErr := imagebytes.ReadU24(reader, imagebytes.LittleEndian)
	return int(widthU24 + 1), int(heightU24 + 1), errors.Join(widthErr, heightErr)
}
