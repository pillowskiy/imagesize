package imagesize_test

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/pillowskiy/imagesize"
)

type TestCase struct {
	Name string
	Path string

	ShouldFail bool
	Expected   *imagesize.ImageInfo
}

type TestGroup struct {
	Name  string
	Cases []TestCase
}

var testCases = []TestGroup{
	{
		Name: "Avif",
		Cases: []TestCase{
			{
				Path: "_testdata/avif/20x20.avif",
				Expected: &imagesize.ImageInfo{
					ImageSize: imagesize.ImageSize{
						Width:  20,
						Height: 20,
					},
					Format: "avif",
				},
			},
		},
	},
	{
		Name: "GIF",
		Cases: []TestCase{
			{
				Path: "_testdata/gif/200x200.gif",
				Expected: &imagesize.ImageInfo{
					ImageSize: imagesize.ImageSize{
						Width:  200,
						Height: 200,
					},
					Format: "gif",
				},
			},
		},
	},
	{
		Name: "HEIC",
		Cases: []TestCase{
			{
				Path: "_testdata/heic/heic.heic",
				Expected: &imagesize.ImageInfo{
					ImageSize: imagesize.ImageSize{
						Width:  2448,
						Height: 3264,
					},
					Format: "heic",
				},
			},

			{
				Name: "MSF1",
				Path: "_testdata/heic/heic_msf1.heic",
				Expected: &imagesize.ImageInfo{
					ImageSize: imagesize.ImageSize{
						Width:  1280,
						Height: 720,
					},
					Format: "heic",
				},
			},
		},
	},
	{
		Name: "JPEG",
		Cases: []TestCase{
			{
				Path: "_testdata/jpeg/20x20.jpg",
				Expected: &imagesize.ImageInfo{
					ImageSize: imagesize.ImageSize{
						Width:  20,
						Height: 20,
					},
					Format: "jpeg",
				},
			},
		},
	},
	{
		Name: "PNG",
		Cases: []TestCase{
			{
				Path: "_testdata/png/20x20.png",
				Expected: &imagesize.ImageInfo{
					ImageSize: imagesize.ImageSize{
						Width:  20,
						Height: 20,
					},
					Format: "png",
				},
			},
			{
				Name: "Animated",
				Path: "_testdata/png/100x100_animated.png",
				Expected: &imagesize.ImageInfo{
					ImageSize: imagesize.ImageSize{
						Width:  100,
						Height: 100,
					},
					Format: "png",
				},
			},
		},
	},
	{
		Name: "WEBP",
		Cases: []TestCase{
			{
				Name: "VP8_",
				Path: "_testdata/webp/vp8_20x20.webp",
				Expected: &imagesize.ImageInfo{
					ImageSize: imagesize.ImageSize{
						Width:  20,
						Height: 20,
					},
					Format: "webp",
				},
			},
			{
				Name: "VP8L",
				Path: "_testdata/webp/vp8l_20x20.webp",
				Expected: &imagesize.ImageInfo{
					ImageSize: imagesize.ImageSize{
						Width:  20,
						Height: 20,
					},
					Format: "webp",
				},
			},
			{
				Name: "VP8X",
				Path: "_testdata/webp/vp8x_180x180.webp",
				Expected: &imagesize.ImageInfo{
					ImageSize: imagesize.ImageSize{
						Width:  180,
						Height: 180,
					},
					Format: "webp",
				},
			},
		},
	},
}

func assertEqual(t *testing.T, expected, actual interface{}, msg string) {
	if expected != actual {
		t.Fatalf("%s: expected %v, got %v", msg, expected, actual)
	}
}

func assertEqualInfo(t *testing.T, expected, actual *imagesize.ImageInfo) {
	assertEqual(t, expected.Format, actual.Format, "Format mismatch")
	assertEqual(t, expected.Width, actual.Width, "Width mismatch")
	assertEqual(t, expected.Height, actual.Height, "Height mismatch")
}

func TestExtractFileInfo(t *testing.T) {
	t.Parallel()

	for _, g := range testCases {
		for _, tt := range g.Cases {
			t.Run(fmt.Sprintf("%s_%s", g.Name, tt.Name), func(t *testing.T) {
				if _, err := os.Stat(tt.Path); os.IsNotExist(err) {
					t.Fatalf("File %s does not exist, nothing to test", tt.Path)
				}

				info, err := imagesize.ExtractFileInfo(tt.Path)
				if err != nil {
					if !tt.ShouldFail {
						t.Fatalf("unexpected error: %v", err)
					}
					if info != nil {
						t.Fatalf("expected nil, got: %v", info)
					}
				} else {
					assertEqualInfo(t, tt.Expected, info)
				}
			})
		}
	}
}

func TestExtractInfo_ReaderPosition_Unchanged(t *testing.T) {
	t.Parallel()

	for _, g := range testCases {
		for _, tt := range g.Cases {
			t.Run(fmt.Sprintf("%s_%s", g.Name, tt.Name), func(t *testing.T) {
				if _, err := os.Stat(tt.Path); os.IsNotExist(err) {
					t.Fatalf("File %s does not exist, nothing to test", tt.Path)
				}

				file, err := os.Open(tt.Path)
				if err != nil {
					t.Fatalf("Failed to open file %s: %v", tt.Path, err)
				}
				defer file.Close()

				initialPos, err := file.Seek(0, io.SeekCurrent)
				if err != nil {
					t.Fatalf("Failed to get initial stream position: %v", err)
				}
				defer func() {
					finalPos, err := file.Seek(0, io.SeekCurrent)
					if err != nil {
						t.Fatalf("Failed to get final stream position: %v", err)
					}

					assertEqual(t, initialPos, finalPos, "Reader should be on initial position")
				}()

				info, err := imagesize.ExtractInfo(file)
				if err != nil {
					if !tt.ShouldFail {
						t.Fatalf("unexpected error: %v", err)
					}
					if info != nil {
						t.Fatalf("expected nil, got: %v", info)
					}
				} else {
					assertEqualInfo(t, tt.Expected, info)
				}
			})
		}
	}
}
