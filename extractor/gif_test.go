package extractor_test

import (
	"bytes"
	"testing"

	"github.com/pillowskiy/imagesize/extractor"
)

func TestGIF(t *testing.T) {
	t.Parallel()
	extractor := extractor.GIF{}

	var (
		gifHeader        = []byte("GIF")
		gifVersionHeader = []byte("89a")
		gifWidth         = []byte{0x01, 0x00} // 1
		gifHeight        = []byte{0x02, 0x00} // 2
	)

	t.Run("BufferSizeMatchesGIFHeaderLength", func(t *testing.T) {
		bufSize := extractor.BufSize()
		expectedBufSize := len(gifHeader)

		if bufSize != expectedBufSize {
			t.Errorf("expected buf size %d, got %d", expectedBufSize, bufSize)
		}
	})

	validGIF := mergeBuffers(
		gifHeader,
		gifVersionHeader,
		gifWidth, gifHeight,
	)

	t.Run("FormatDetection", func(t *testing.T) {
		format, matched := extractor.MatchFormat(validGIF)
		if !matched {
			t.Error("expected match for valid GIF file")
		}

		expectedFormat := "gif"
		if format != expectedFormat {
			t.Errorf("expected format %s, got %s", expectedFormat, format)
		}
	})

	t.Run("ExtractSizeFromValidImage", func(t *testing.T) {
		reader := bytes.NewReader(validGIF)
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

	t.Run("CorruptedImage", func(t *testing.T) {
		invalidGIF := mergeBuffers(
			gifHeader,
			gifVersionHeader,
			gifWidth,
		)

		reader := bytes.NewReader(invalidGIF)
		_, _, err := extractor.ExtractSize(reader)

		if err == nil {
			t.Fatalf("expected error due to missing height, got nil")
		}
	})

	t.Run("BufferSizeMatchesJPEGHeaderLength", func(t *testing.T) {
		nonGIF := []byte("NOTGIFHEADER")
		_, matched := extractor.MatchFormat(nonGIF)

		if matched {
			t.Error("expected no match for non-GIF file")
		}
	})
}
