package imagesize_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/pillowskiy/imagesize"
)

func BenchmarkSize(b *testing.B) {
	imagesDir := "_testdata"
	var paths []string

	err := filepath.Walk(imagesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			absPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			paths = append(paths, absPath)
		}
		return nil
	})
	if err != nil {
		b.Fatalf("Failed to walk the directory: %v", err)
	}

	b.ResetTimer()
	for _, path := range paths {
		b.Run(fmt.Sprintf("Info_%s", path), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if _, err := imagesize.ExtractFileInfo(path); err != nil {
					b.Fatalf("Failed to decode image: %v", err)
				}
			}
		})
	}
}
