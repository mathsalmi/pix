package main

import (
	"image"
)

// Image is an extension of the image.Image type which
// contains some methods for convenience
type Image struct {
	image.Image
}

// Width returns the image width
func (i Image) Width() int {
	b := i.Bounds()
	return b.Max.X - b.Min.X
}

// Height returns the image height
func (i Image) Height() int {
	b := i.Bounds()
	return b.Max.Y - b.Min.Y
}
