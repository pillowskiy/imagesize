package extractor_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/pillowskiy/imagesize/extractor"
)

func TestJPEG(t *testing.T) {
	t.Parallel()
	extractor := extractor.JPEG{}

	jpegHeader := []byte{0xFF, 0xD8, 0xFF, 0xE0} // SOI and APP0 markers

	validJPEG := mergeBuffers(
		jpegHeader,
		[]byte{
			0x00, 0x13, // Length of the marker (17 bytes - until 0xFF (another marker))
			0x4A, 0x46, 0x49, 0x46, 0x00, 0x01, // "JFIF" identifier
			0x01, 0x01, // Version
			0x00,       // Units
			0x00, 0x01, // X density
			0x00, 0x01, // Y density
			0x00, 0x00, // X thumbnail
			0x00, 0x00, // Y thumbnail
			0xFF, 0xC0, // SOF0 marker
			0x00, 0x0B, // Length of the marker
			0x08,       // Precision (8 bits)
			0x00, 0x02, // Height
			0x00, 0x01, // Width
			0x01,       // Number of components
			0x01,       // Component ID
			0x00,       // Horizontal sampling factor
			0x00,       // Vertical sampling factor
			0xFF, 0xDA, // SOS marker
			// 0x00, 0x0C, // Length of the marker
			// 0x01,             // Number of components in the scan
			// 0x01,             // Component ID
			// 0x00,             // Spectral selection
			// 0x00,             // Successive approximation
			// 0x00, 0x00, 0x00, // Compressed data
			0xFF, 0xD9, // EOI (End of Image)
		},
	)

	t.Run("buf size should be a length of the jpeg header", func(t *testing.T) {
		bufSize := extractor.BufSize()
		expectedBufSize := len(jpegHeader)

		if bufSize != expectedBufSize {
			t.Errorf("expected buf size %d, got %d", expectedBufSize, bufSize)
		}
	})

	t.Run("should match format", func(t *testing.T) {
		format, matched := extractor.MatchFormat(validJPEG)
		if !matched {
			t.Error("expected match for valid JPEG file")
		}

		expectedFormat := "jpeg"
		if format != expectedFormat {
			t.Errorf("expected format %s, got %s", expectedFormat, format)
		}
	})

	t.Run("should extract size", func(t *testing.T) {
		reader := bytes.NewReader(validJPEG)
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

	t.Run("should error on corrupted JPEG", func(t *testing.T) {
		invalidJPEG := mergeBuffers(
			jpegHeader,
			[]byte{
				0x00, 0x12, // Incorrect length of the marker 16 instead of 17
				0x4A, 0x46, 0x49, 0x46, 0x00, 0x01, // "JFIF" identifier
				0x01, 0x01, // Version
				0x00,       // Units
				0x00, 0x01, // X density
				0x00, 0x01, // Y density
				0x00, 0x00, // X thumbnail
				0x00, 0x00, // Y thumbnail
				0xFF, 0xC0, // SOF0 marker
				0x00, 0x0B, // Length of the marker
				0x08,       // Precision (8 bits)
				0x00, 0x02, // Height
				0x00, 0x01, // Width
				0x01,       // Number of components
				0x01,       // Component ID
				0x00,       // Horizontal sampling factor
				0x00,       // Vertical sampling factor
				0xFF, 0xDA, // SOS marker
				0xFF, 0xD9, // EOI (End of Image)
			},
		)

		reader := bytes.NewReader(invalidJPEG)
		_, _, err := extractor.ExtractSize(reader)

		if err == nil {
			t.Fatalf("expected error due to missing height, got nil")
		}
	})

	t.Run("should EOF on bad JPEG", func(t *testing.T) {
		invalidJPEG := mergeBuffers(
			jpegHeader,
			[]byte{
				0x00, 0x13, // Length of the marker (17 bytes - until 0xFF (another marker))
				0x4A, 0x46, 0x49, 0x46, 0x00, 0x01, // "JFIF" identifier
				0x01, 0x01, // Version
				0x00,       // Units
				0x00, 0x01, // X density
				0x00, 0x01, // Y density
				0x00, 0x00, // X thumbnail
				0x00, 0x00, // Y thumbnail
				0xFF, 0xC0, // SOF0 marker
				0x00, 0x0B, // Length of the marker
				0x08,       // Precision (8 bits)
				0x00, 0x02, // Height
			},
		)

		reader := bytes.NewReader(invalidJPEG)
		_, _, err := extractor.ExtractSize(reader)

		if err == nil {
			t.Fatalf("expected error due to missing height, got nil")
		}

		if !errors.Is(err, io.EOF) {
			t.Fatalf("expected EOF error")
		}
	})

	t.Run("should not match format for non-JPEG file", func(t *testing.T) {
		nonJPEG := []byte("NOTJPEGHEADER")
		_, matched := extractor.MatchFormat(nonJPEG)

		if matched {
			t.Error("expected no match for non-JPEG file")
		}
	})
}
