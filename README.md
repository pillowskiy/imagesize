[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat)](https://pkg.go.dev/github.com/pillowskiy/imagesize)
[![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/pillowskiy/imagesize/ci.yml?style=flat)](https://github.com/pillowskiy/imagesize/actions)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat)](https://raw.githubusercontent.com/pillowskiy/imagesize/master/LICENSE)

# GO Imagesize 
Efficiently extract the dimensions of various image formats without fully loading the files into memory.

This library is designed to minimize overhead when determining the fomat, width and height of supported image types. By reading only the necessary portions of the file, `imagesize` is lightweight, fast, and has minimal dependencies.

----

## Getting Started

Install the package in your Go project:
```bash
go get github.com/pillowskiy/imagesize@latest
```

## Supported Formats

The library currently supports the following image formats:
- avif
- gif
- heic / heif
- jpeg
- png
- webp

If you need support for additional formats, feel free to open an issue or contribute!

## Examples

### Extracting Image Information from a File
```go
package main

import (
	"fmt"
	"log"
	"github.com/pillowskiy/imagesize"
)

func main() {
	// Path to the image file
	imagePath := "test.png"

	info, err := imagesize.ExtractFileInfo(imagePath)
	if err != nil {
		log.Fatalf("Error extracting image info: %v", err)
	}
	fmt.Printf("Result: %+v\n", info)
}
```

### Extracting Image Information from a Byte Slice (Blob)

```go
package main

import (
	"fmt"
	"log"
	"github.com/pillowskiy/imagesize"
)

func main() {
	// Example image data (first few bytes of a PNG image)
	data := []byte{
		0x89, 0x50, 0x4E, 0x47, // PNG Header
		0x0D, 0x0A, 0x1A, 0x0A, // Sequence
		0x00, 0x00, 0x00, 0x0D, // IHDR Length
		0x49, 0x48, 0x44, 0x52, // IHDR Header
		0x00, 0x00, 0x00, 0x01, // 1 width in big-endian
		0x00, 0x00, 0x00, 0x02, // 2 height in big-endian
		// Additional image data...
	}

	// Extract image info from the byte slice
	info, err := imagesize.ExtractBlobInfo(data)
	if err != nil {
		log.Fatalf("Error extracting image info: %v", err)
	}
	fmt.Printf("Result: %+v\n", info)
}
```

## Inspiration

While working on my side project, I found that getting basic image information usually means decoding the whole image. I couldn't find a suitable Go library for this, but I found a similar library in Rust ([Roughsketch/imagesize](https://github.com/Roughsketch/imagesize)). I didn't want to set up an RPC service or use WASM, so I decided to create my own solution.