# goincv

Goincv is a collection of Go libraries for scientific computing and image processing. The libraries include tools for manipulating and analyzing images, as well as common utilities for scientific computing tasks.

## Installation

To use goincv in your Go project, simply run:

```
go get github.com/wailovet/goincv
```

## Usage

To use goincv in your Go code, import the package with:

```
import "github.com/wailovet/goincv"
```

### Image processing

Goincv provides several functions for manipulating and analyzing images, including:

- `goincv.ImRead()`: Convert an image to a [][][]uint8 of pixel values 

## Examples

Here is an example of how to use goincv to resize an image:

```go
package main

import (
	"fmt"
	"image"
	"github.com/wailovet/goincv"
)

func main() {
	// Open the image file
	img := goincv.File2Image(imgPath)
	timg := resize.Resize(uint(128), uint(128), img, resize.Bilinear)
	data := goincv.ImageMeanNormalizeF32RGBCHW(goincv.ImReadRGBCHW(timg), []float32{0.485, 0.456, 0.406}, []float32{0.229, 0.224, 0.225})
	fData := goincv.FlattenFloat32(data)
	fmt.Println(fData)
}

```
 