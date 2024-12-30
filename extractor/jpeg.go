package extractor

import (
	"bytes"
	"fmt"
	"io"
)

var jpegHeader = []byte("\xFF\xD8\xFF")

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
		err = fmt.Errorf("failed to read byte: %w", err)
		return
	}

	for buf1[0] != 0xDA && buf1[0] != 0 {
		// Read until 0xFF
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

		// Check for specific markers (0xC0 to 0xC3)
		if buf1[0] >= 0xC0 && buf1[0] <= 0xC3 {
			_, err = reader.Seek(3, io.SeekCurrent)
			if err != nil {
				err = fmt.Errorf("failed to seek 3 bytes forward: %w", err)
				return
			}

			// Read the width and height (same as structure!(">HH") in Rust)
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
			err = fmt.Errorf("failed to read 2 bytes: %w", err)
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
