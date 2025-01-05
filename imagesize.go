package imagesize

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/pillowskiy/imagesize/imagerrors"
)

// Extracts image info from a file at the specified path.
// It opens the file, reads its contents, and then delegates to ExtractInfo to parse the metadata.
func ExtractFileInfo(path string) (*ImageInfo, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	return ExtractInfo(file)
}

// Extracts image info from a byte slice.
// It creates a bytes.Reader for the byte slice and delegates to ExtractInfo.
func ExtractBlobInfo(buf []byte) (*ImageInfo, error) {
	reader := bytes.NewReader(buf)
	return ExtractInfo(reader)
}

// Extracts image info from an io.ReaderAt.
// It ensures the provided reader is compatible with the underlying extraction logic,
// converting it to an io.ReadSeeker if necessary, and then delegates to extractInfo.
func ExtractInfo(reader io.ReaderAt) (*ImageInfo, error) {
	var sr io.ReadSeeker

	if r, ok := reader.(io.ReadSeeker); ok {
		sr = r
	} else {
		const maxInt int64 = 1<<63 - 1
		sr = io.NewSectionReader(reader, 0, maxInt)
	}

	return extractInfo(sr)
}

func extractInfo(reader io.ReadSeeker) (info *ImageInfo, err error) {
	buf, err := readAtLeast(reader, make([]byte, 0, 8), 4)
	if err != nil {
		return nil, err
	}

	defer func() {
		if _, seekErr := reader.Seek(0, io.SeekStart); seekErr != nil {
			err = imagerrors.Join(err, seekErr)
		}
	}()

	info = new(ImageInfo)
	for _, ext := range imageSizeExtractors {
		reqBuf := ext.BufSize()
		if len(buf) <= reqBuf {
			buf, err = readAtLeast(reader, buf, reqBuf)
			if err != nil {
				return nil, err
			}
		}

		format, match := ext.MatchFormat(buf)
		if !match {
			continue
		}

		info.Format = format

		width, height, err := ext.ExtractSize(reader)
		if err != nil {
			return nil, err
		}

		info.Width = width
		info.Height = height

		return info, err
	}

	return nil, errors.New("unknown format")
}

// Ensures the buffer contains at least the required number of bytes (`needed`).
// If the buffer's length is already sufficient, it returns the buffer unchanged.
// Otherwise, it reads additional bytes from the provided io.Reader to meet the requirement.
func readAtLeast(reader io.Reader, buf []byte, needed int) ([]byte, error) {
	if len(buf) >= needed {
		return buf, nil
	}

	extra := make([]byte, needed-len(buf))
	_, err := reader.Read(extra)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return append(buf, extra...), nil
}
