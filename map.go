package imagesize

import "github.com/pillowskiy/imagesize/extractor"

var imageSizeExtractors = []SizeExtractor{
	extractor.JPEG{},
	extractor.GIF{},
	extractor.WEBP{},
	extractor.PNG{},
}

// RegisterSizeExtractor adds a new SizeExtractor to the list of image size extractors.
// This allows dynamic extension of supported formats without modifying the original slice.
func RegisterSizeExtractor(e SizeExtractor) {
	imageSizeExtractors = append(imageSizeExtractors, e)
}
