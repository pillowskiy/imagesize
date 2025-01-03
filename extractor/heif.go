package extractor

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/pillowskiy/imagesize/imagebytes"
)

var ftypHeader = []byte("ftyp")

const (
	hevcBrandKey uint8 = 1
	av1BrandKey  uint8 = 2
	jpegBrandKey uint8 = 3
)

var ftypCompatibleBrandsMap = map[string]struct{}{
	"mif1": {},
	"msf1": {},
	"mif2": {},
	"miaf": {},
}

var brandsKeyMap = map[string]uint8{
	"heic": hevcBrandKey,
	"heix": hevcBrandKey,
	"heis": hevcBrandKey,
	"hevs": hevcBrandKey,
	"heim": hevcBrandKey,
	"hevm": hevcBrandKey,
	"hevc": hevcBrandKey,
	"hevx": hevcBrandKey,

	"avif": av1BrandKey,
	"avio": av1BrandKey,
	"avis": av1BrandKey,
	"MA1A": av1BrandKey,
	"MA1B": av1BrandKey,

	"jpeg": jpegBrandKey,
	"jpgs": jpegBrandKey,
}

// See:
//   - https://en.wikipedia.org/wiki/High_Efficiency_Image_File_Format
//   - https://github.com/nokiatech/heif
type HEIF struct{}

func (e HEIF) ExtractSize(reader io.ReadSeeker) (width, height int, err error) {
	if _, err = reader.Seek(0, io.SeekStart); err != nil {
		return
	}

	// Read the ftyp header size
	ftypSize, err := imagebytes.ReadU32(reader, imagebytes.BigEndian)
	if err != nil {
		err = fmt.Errorf("failed to read ftyp header size: %w", err)
		return
	}

	// Jump to the first actual box offset
	if _, err = reader.Seek(int64(ftypSize), io.SeekStart); err != nil {
		return
	}

	// Skip to meta tag
	if _, err = e.skipToTag(reader, []byte("meta")); err != nil {
		return
	}

	// Discard the junk value after meta tag
	if _, err = reader.Seek(4, io.SeekCurrent); err != nil {
		return
	}

	// Skip to iprp tag
	if _, err = e.skipToTag(reader, []byte("iprp")); err != nil {
		return
	}

	// Find ipco tag
	ipcoSizeU32, err := e.skipToTag(reader, []byte("ipco"))
	if err != nil {
		return
	}
	ipcoSize := int(ipcoSizeU32)

	// Keep track of the max size of ipco tag
	var maxWidth, maxHeight int
	foundIspe := false
	var rotation uint8

	// Loop through all tags within ipco
	for {
		tag, size, tagErr := imagebytes.ReadTag(reader)
		if tagErr != nil {
			err = tagErr
			return
		}

		// Size of tag length + tag cannot be under 8 (4 bytes each)
		if size < 8 {
			err = errors.New("corrupted image")
			return
		}

		// Image spatial Extents (ispe) indicates the width and height of the associated image item
		if tag == "ispe" {
			foundIspe = true
			// Discard junk value
			if _, seekErr := reader.Seek(4, io.SeekCurrent); seekErr != nil {
				err = seekErr
				return
			}

			widthU32, widthErr := imagebytes.ReadU32(reader, imagebytes.BigEndian)
			heightU32, heightErr := imagebytes.ReadU32(reader, imagebytes.BigEndian)
			if readSizeErr := errors.Join(widthErr, heightErr); err != nil {
				err = readSizeErr
				return
			}

			w, h := int(widthU32), int(heightU32)
			// Assign new largest size by area
			if w*h > maxWidth*maxHeight {
				maxWidth = w
				maxHeight = h
			}
		} else if tag == "irot" { // Image Rotation (irot): Rotation by 90, 180, or 270 degrees.
			rotation, err = imagebytes.ReadU8(reader)
			if err != nil {
				return
			}
		} else if size >= ipcoSize { // If we've gone past the ipco boundary, then break
			break
		} else {
			// If we're still inside ipco, consume all bytes for the current tag
			ipcoSize -= size

			if _, seekErr := reader.Seek(int64(size-8), io.SeekCurrent); seekErr != nil {
				err = seekErr
				return
			}
		}
	}

	if !foundIspe {
		err = errors.New("not enough data to extract size: ispe not found")
		return
	}

	// If rotation is 90deg (1) or 270deg (3), swap dims
	if rotation == 1 || rotation == 3 {
		maxWidth, maxHeight = maxHeight, maxWidth
	}

	width = maxWidth
	height = maxHeight

	return
}

func (e HEIF) BufSize() int {
	return 24
}

func (e HEIF) MatchFormat(buf []byte) (format string, match bool) {
	if len(buf) < 12 || !bytes.Equal(buf[4:8], ftypHeader) {
		return
	}

	var brand [4]byte
	copy(brand[:], buf[8:12])

	if ext, match := e.matchBrandFormat(brand); match {
		return ext, match
	}

	return e.matchCompatibleBrandFormat(buf, 12, 24)
}

// See: https://github.com/nokiatech/heif/blob/be43efdf273ae9cf90e552b99f16ac43983f3d19/srcs/reader/heifreaderimpl.cpp#L738
func (e HEIF) matchCompatibleBrandFormat(buf []byte, startIndex, endIndex int) (string, bool) {
	size := len(buf)
	if size < endIndex {
		endIndex = size - size%4
	}

	const chunkSize = 4

	for i := startIndex; i < endIndex; i += chunkSize {
		eob := i + chunkSize // end of potential brand chunk index
		b := buf[i:eob]
		if _, ok := ftypCompatibleBrandsMap[string(b)]; ok {
			var nextBrand [chunkSize]byte
			copy(nextBrand[:], buf[eob:eob+chunkSize])

			if format, match := e.matchBrandFormat(nextBrand); match {
				return format, match
			}
		}
	}

	return "", false
}

func (e HEIF) matchBrandFormat(brand [4]byte) (format string, match bool) {
	matchKey, match := brandsKeyMap[string(brand[:])]
	if !match {
		return "", false
	}

	switch matchKey {
	case hevcBrandKey:
		format = "heic"
	case av1BrandKey:
		format = "avif"
	case jpegBrandKey:
		format = "jpeg"
	}

	return
}

func (e HEIF) skipToTag(reader io.ReadSeeker, tag []byte) (uint32, error) {
	var tagBuf [4]byte

	for {
		size, err := imagebytes.ReadU32(reader, imagebytes.BigEndian)
		if err != nil {
			return 0, err
		}

		if _, err := io.ReadFull(reader, tagBuf[:]); err != nil {
			return 0, err
		}

		if bytes.Equal(tagBuf[:], tag) {
			return size, nil
		}

		if size >= 8 {
			if _, err := reader.Seek(int64(size)-8, io.SeekCurrent); err != nil {
				return 0, err
			}
		} else {
			return 0, errors.New("invalid HEIF box size")
		}
	}
}
