package main

import (
	"image"

	"github.com/anthonynsimon/bild/transform"
)

func applyResize(image *image.Image, w, h int) error {
	*image = transform.Resize(*image, w, h, transform.Linear)
	return nil
}
