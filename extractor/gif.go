package extractor

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

var gifHeader = []byte("\x47\x49\x46")

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

	widthU16, widthErr := readU16(reader, LittleEndian)
	heightU16, heightErr := readU16(reader, LittleEndian)
	return int(widthU16), int(heightU16), errors.Join(widthErr, heightErr)
}
