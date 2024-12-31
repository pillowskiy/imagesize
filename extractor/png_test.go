package extractor_test

import (
	"bytes"
	"testing"

	"github.com/pillowskiy/imagesize/extractor"
)

func TestPNG(t *testing.T) {
	t.Parallel()
	extractor := extractor.PNG{}

	var (
		pngHeader         = []byte("\x89\x50\x4E\x47")     // PNG header
		pngSequenceHeader = []byte{0x0D, 0x0A, 0x1A, 0x0A} // special sequence for PNG
		ihdrLengthHeader  = []byte{0x00, 0x00, 0x00, 0x0D} // IHDR length (13 bytes)
		ihdrHeader        = []byte("IHDR")                 // IHDR chunk header
		pngWidth          = []byte{0x00, 0x00, 0x00, 0x01} // 1 (width in big-endian)
		pngHeight         = []byte{0x00, 0x00, 0x00, 0x02} // 2 (height in big-endian)
	)

	mergePNG := func() []byte {
		return mergeBuffers(
			pngHeader,
			pngSequenceHeader,
			ihdrLengthHeader,
			ihdrHeader,
			pngWidth, pngHeight,
		)
	}

	t.Run("buf size should match PNG header length", func(t *testing.T) {
		bufSize := extractor.BufSize()
		expectedBufSize := len(pngHeader)

		if bufSize != expectedBufSize {
			t.Errorf("expected buf size %d, got %d", expectedBufSize, bufSize)
		}
	})

	t.Run("should match format for valid PNG", func(t *testing.T) {
		validPNG := mergePNG()
		format, matched := extractor.MatchFormat(validPNG)
		if !matched {
			t.Error("expected match for valid PNG file")
		}

		expectedFormat := "png"
		if format != expectedFormat {
			t.Errorf("expected format %s, got %s", expectedFormat, format)
		}
	})

	t.Run("should extract size for valid PNG", func(t *testing.T) {
		validPNG := mergePNG()
		reader := bytes.NewReader(validPNG)
		width, height, err := extractor.ExtractSize(reader)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if width != 1 {
			t.Errorf("expected width 1, got %d", width)
		}

		if height != 2 {
			t.Errorf("expected height 2, got %d", height)
		}
	})

	t.Run("should error on corrupted PNG file (missing IHDR)", func(t *testing.T) {
		invalidPNG := mergeBuffers(
			pngHeader,
			pngSequenceHeader,
			ihdrHeader,
			// Missing IHDR header and image data
		)

		reader := bytes.NewReader(invalidPNG)
		_, _, err := extractor.ExtractSize(reader)

		if err == nil {
			t.Fatalf("expected error due to missing IHDR header, got nil")
		}
	})

	t.Run("should not match format for non-PNG file", func(t *testing.T) {
		nonPNG := []byte("NOTPNGHEADER")
		_, matched := extractor.MatchFormat(nonPNG)

		if matched {
			t.Error("expected no match for non-PNG file")
		}
	})
}
