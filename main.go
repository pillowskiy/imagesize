package imagesize

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
)

func ExtractFileInfo(path string) (*ImageInfo, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	return ExtractInfo(file)
}

func ExtractBlobInfo(buf []byte) (*ImageInfo, error) {
	reader := bytes.NewReader(buf)
	return ExtractInfo(reader)
}

func ExtractInfo(reader io.ReadSeeker) (*ImageInfo, error) {
	return extractInfo(reader)
}

func extractInfo(reader io.ReadSeeker) (info *ImageInfo, err error) {
	buf, err := readAtLeast(reader, make([]byte, 0, 2), 2)
	if err != nil {
		return nil, err
	}

	defer func() {
		if _, seekErr := reader.Seek(0, io.SeekStart); seekErr != nil {
			err = errors.Join(err, seekErr)
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
