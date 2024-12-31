package extractor_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/pillowskiy/imagesize/extractor"
)

func TestWEBP(t *testing.T) {
	t.Parallel()
	extractor := extractor.WEBP{}

	var (
		webpRIFFHeader     = []byte("RIFF")
		webpFileSizeHeader = []byte{0x00, 0x00, 0x00, 0x10}
		webpFormatHeader   = []byte("WEBP")
	)

	validWEBP := mergeBuffers(
		webpRIFFHeader,
		webpFileSizeHeader,
		webpFormatHeader,
	)

	validWEBPs := []struct {
		Name string
		Buf  []byte
	}{
		{
			Name: "VP8 ",
			Buf: mergeBuffers(
				validWEBP,
				[]byte("VP8 "),
				make([]byte, 10),   // See https://datatracker.ietf.org/doc/html/rfc6386#section-1
				[]byte{0x01, 0x00}, // Width: 1 (u16 little endian)
				[]byte{0x02, 0x00}, // Width: 2 (u16 little endian)
			),
		},
		{
			Name: "VP8X",
			Buf: mergeBuffers(
				validWEBP,
				[]byte("VP8X"),
				// Reserved, ICC, Alpha, Exif metadata, Animation...
				make([]byte, 8),          // See https://developers.google.com/speed/webp/docs/riff_container#extended_header
				[]byte{0x01, 0x00, 0x00}, // Width: 1 (u24 little endian)
				[]byte{0x02, 0x00, 0x00}, // Width: 2 (u24 little endian)
			),
		},
		{
			Name: "VP8L",
			Buf: mergeBuffers(
				validWEBP,
				[]byte("VP8L"),
				// See https://developers.google.com/speed/webp/docs/webp_lossless_bitstream_specification#3_riff_header
				make([]byte, 5),
				// 14 bits for width and 14 bits for height:
				[]byte{0x01, 0x80, 0x00},       // Width: 1, Height: 2 as 0x01 | (0x02 << 14)
				[]byte{0x00, 0x00, 0x00, 0x00}, // alpha_is_used and version_number (3 bit code that must be set to 0)
			),
		},
	}

	for _, validWEBP := range validWEBPs {
		t.Run(fmt.Sprintf("should extract correct size for %s WebP", validWEBP.Name), func(t *testing.T) {
			reader := bytes.NewReader(validWEBP.Buf)
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
	}

	t.Run("should error on invalid VP8 WebP format", func(t *testing.T) {
		invalidWEBP := mergeBuffers(
			webpRIFFHeader,
			webpFileSizeHeader,
			webpFormatHeader,
			[]byte("VP8Y"),
		)

		reader := bytes.NewReader(invalidWEBP)
		_, _, err := extractor.ExtractSize(reader)

		if err == nil {
			t.Fatalf("expected error due to invalid VP8 format, got nil")
		}
	})

	t.Run("should not match format for non-WebP file", func(t *testing.T) {
		_, matched := extractor.MatchFormat([]byte("NOTWEBPHEADER"))
		if matched {
			t.Error("expected no match for non-WebP file")
		}
	})
}
