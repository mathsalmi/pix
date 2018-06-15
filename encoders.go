package main

import (
	"image"
	"image/gif"
	"io"

	"github.com/anthonynsimon/bild/imgio"
)

// GIFEncoder returns an encoder to GIF
func GIFEncoder(NumColors int) imgio.Encoder {
	return func(w io.Writer, img image.Image) error {
		return gif.Encode(w, img, &gif.Options{NumColors: NumColors})
	}
}
