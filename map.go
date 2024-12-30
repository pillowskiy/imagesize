package imagesize

import "github.com/pillowskiy/imagesize/extractor"

var imageSizeExtractors = []SizeExtractor{
	extractor.JPEG{},
	extractor.GIF{},
	extractor.WEBP{},
	extractor.PNG{},
}

func RegisterSizeExtractor(e SizeExtractor) {
	imageSizeExtractors = append(imageSizeExtractors, e)
}
