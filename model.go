package imagesize

import "io"

// Interface for extracting size from various formats.
type SizeExtractor interface {
	// BufSize returns the size of the buffer needed to identify the format.
	// This method helps ensure the caller knows how many bytes to read for format detection.
	BufSize() int

	// MatchFormat checks if the provided buffer matches the format associated with the implementation.
	// Parameters:
	//   - buf: A byte slice containing the data to inspect.
	// Returns:
	//   - A string representing the detected format (e.g., "JPEG", "PNG").
	//   - A boolean indicating whether the buffer matches the format.
	MatchFormat(buf []byte) (string, bool)

	// ExtractSize extracts the dimensions (width and height) of an object from the provided reader.
	// Parameters:
	//   - reader: An io.ReadSeeker __positioned at the beginning__ of the object to analyze.
	// Returns:
	//   - width: The width of the object.
	//   - height: The height of the object.
	//   - err: An error if the dimensions cannot be extracted.
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
