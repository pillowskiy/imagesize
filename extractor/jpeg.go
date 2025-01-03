package extractor

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/pillowskiy/imagesize/imagebytes"
)

// Without APP marker
// See https://stackoverflow.com/questions/5413022/is-the-2nd-and-3rd-byte-of-a-jpeg-image-always-the-app0-or-app1-marker
var jpegHeader = []byte("\xFF\xD8\xFF")

// JPEG defines an extractor for JPEG image format.
//
// The JPEG file format starts with a specific sequence of bytes for its header:
//
// 1. The first 2 bytes contain the value 0xFF, 0xD8, which represent the Start of Image (SOI) marker.
// 2. The following bytes represent various segments within the JPEG file.
// Each segment begins with a 2-byte marker that starts with 0xFF (indicating the start of a new segment)
// 3. After the SOI marker, there is typically a segment containing the JPEG quantization table (0xDB),
// followed by other segments, including the Huffman table (0xC4), Start of Frame (0xC0), and the Start of Scan (0xDA) markers, which contain the image data.
// 4. The file ends with a 2-byte marker (0xFF, 0xD9), known as the End of Image (EOI) marker. This marks the conclusion of the JPEG file.
type JPEG struct{}

func (e JPEG) BufSize() int {
	return len(jpegHeader)
}

func (e JPEG) MatchFormat(buf []byte) (string, bool) {
	return "jpeg", bytes.HasPrefix(buf, jpegHeader)
}

func (e JPEG) ExtractSize(reader io.ReadSeeker) (width, height int, err error) {
	if _, err = reader.Seek(2, io.SeekStart); err != nil {
		err = fmt.Errorf("failed to seek to position: %w", err)
		return
	}

	var seg [1]byte
	if seg, err = e.readSegment(reader); err != nil {
		return
	}

	// Read untill Start of Scan marker or end of file
	for seg[0] != 0xDA && seg[0] != 0 {
		// Read until 0xFF (Start of Segment)
		for seg[0] != 0xFF {
			if seg, err = e.readSegment(reader); err != nil {
				return
			}
		}

		// Skip past all 0xFF bytes
		for seg[0] == 0xFF {
			if seg, err = e.readSegment(reader); err != nil {
				return
			}
		}

		// Check for specific markers (0xC0 to 0xC3 - Start of Frame, contains image size data)
		if seg[0] >= 0xC0 && seg[0] <= 0xC3 {
			_, err = reader.Seek(3, io.SeekCurrent)
			if err != nil {
				return
			}

			heightU16, heightErr := imagebytes.ReadU16(reader, imagebytes.BigEndian)
			widthU16, widthErr := imagebytes.ReadU16(reader, imagebytes.BigEndian)

			if sizeErr := errors.Join(widthErr, heightErr); sizeErr != nil {
				err = fmt.Errorf("failed to read image size: %w", sizeErr)
				return
			}

			return int(widthU16), int(heightU16), nil
		}

		offset, offsetErr := imagebytes.ReadU16(reader, imagebytes.BigEndian)
		if err != nil {
			err = fmt.Errorf("failed to read offset: %w", offsetErr)
			return
		}

		if _, err = reader.Seek(int64(offset-2), io.SeekCurrent); err != nil {
			err = fmt.Errorf("failed to seek to the next segment: %w", err)
			return
		}

		if seg, err = e.readSegment(reader); err != nil {
			return
		}
	}

	err = errors.New("failed to read image size, stop marker was reached")
	return
}

// More readable, go compiler should inline this, so no performance loss
func (e JPEG) readSegment(reader io.ReadSeeker) (seg [1]byte, err error) {
	if _, err = reader.Read(seg[:]); err != nil {
		err = fmt.Errorf("failed to read segment: %w", err)
	}
	return
}
