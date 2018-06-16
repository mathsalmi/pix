package main

import (
	"image"

	"github.com/anthonynsimon/bild/transform"
)

// ApplyEffects applies effects and transformations to the given image
func ApplyEffects(img *image.Image, options options) {
	applyResize(img, options)
}

// applyResize resizes the given image given the options
func applyResize(image *image.Image, options options) error {
	width, height, err := options.Resize()
	if err != nil {
		return ErrEffectNotApplied
	}

	*image = transform.Resize(*image, width, height, transform.Linear)
	return nil
}
