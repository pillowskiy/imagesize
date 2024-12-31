package extractor

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/pillowskiy/imagesize/imagebytes"
)

var gifHeader = []byte("\x47\x49\x46")

// GIF defines an extractor for the GIF image format.
//
// The GIF file format starts with a fixed header structure:
// 1. The first 3 bytes contain the ASCII characters "GIF", identifying the file as a GIF image.
// 2. The next 3 bytes specify the version of the GIF format. ("89a", "87a" etc.)
//
// After the header, the GIF file contains the Logical Screen Descriptor (LSD), which specifies the width and height of the image in pixels:
// 3. The next 2 bytes represent the width of the image (unsigned 16-bit integer, little-endian).
// 4. The following 2 bytes represent the height of the image (unsigned 16-bit integer, little-endian).
type GIF struct{}

func (e GIF) BufSize() int {
	return len(gifHeader)
}

func (e GIF) MatchFormat(buf []byte) (string, bool) {
	return "gif", bytes.HasPrefix(buf, gifHeader)
}

func (e GIF) ExtractSize(reader io.ReadSeeker) (width, height int, err error) {
	if _, err = reader.Seek(6, io.SeekStart); err != nil {
		err = fmt.Errorf("failed to seek: %w", err)
		return
	}

	widthU16, widthErr := imagebytes.ReadU16(reader, imagebytes.LittleEndian)
	heightU16, heightErr := imagebytes.ReadU16(reader, imagebytes.LittleEndian)
	return int(widthU16), int(heightU16), errors.Join(widthErr, heightErr)
}
