package imagesize

import "io"

type SizeExtractor interface {
	BufSize() int
	MatchFormat(buf []byte) (string, bool)
	ExtractSize(reader io.ReadSeeker) (width int, height int, err error)
}

type ImageSize struct {
	Width  int
	Height int
}

type ImageInfo struct {
	ImageSize
	Format string
}
