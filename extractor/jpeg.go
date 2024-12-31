package extractor

import (
	"bytes"
	"fmt"
	"io"
)

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

	buf1 := make([]byte, 1)
	buf2 := make([]byte, 2)
	buf4 := make([]byte, 4)

	if _, err = reader.Read(buf1); err != nil {
		return
	}

	// Read untill Start of Scan marker or end of file
	for buf1[0] != 0xDA && buf1[0] != 0 {
		// Read until 0xFF (Start of Segment)
		for buf1[0] != 0xFF {
			_, err = reader.Read(buf1)
			if err != nil {
				err = fmt.Errorf("failed to read byte: %w", err)
				return
			}
		}

		// Skip past all 0xFF bytes
		for buf1[0] == 0xFF {
			_, err = reader.Read(buf1)
			if err != nil {
				err = fmt.Errorf("failed to read byte: %w", err)
				return
			}
		}

		// Check for specific markers (0xC0 to 0xC3 - Start of Frame, contains image size data)
		if buf1[0] >= 0xC0 && buf1[0] <= 0xC3 {
			_, err = reader.Seek(3, io.SeekCurrent)
			if err != nil {
				return
			}

			_, err = reader.Read(buf4)
			if err != nil {
				err = fmt.Errorf("failed to read width and height: %w", err)
				return
			}

			// Unpack the width and height (big-endian)
			width = int(buf4[0])<<8 | int(buf4[1])
			height = int(buf4[2])<<8 | int(buf4[3])

			break
		}

		if _, err = reader.Read(buf2); err != nil {
			return
		}

		// Unpack the offset (big-endian)
		offset := int(buf2[0])<<8 | int(buf2[1])
		if _, err = reader.Seek(int64(offset-2), io.SeekCurrent); err != nil {
			err = fmt.Errorf("failed to seek to the next segment: %w", err)
			return
		}

		if _, err = reader.Read(buf1); err != nil {
			err = fmt.Errorf("failed to read byte: %w", err)
			return
		}
	}

	return
}
