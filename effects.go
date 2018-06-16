package main

import (
	"image"

	"github.com/anthonynsimon/bild/transform"
)

// ApplyEffects applies effects and transformations to the given image
func ApplyEffects(img *image.Image, options options) {
	applyCrop(img, options)
	applyResize(img, options)
}

// applyResize resizes the given image given the options
func applyResize(img *image.Image, options options) error {
	width, height, err := options.Resize()
	if err != nil {
		return ErrEffectNotApplied
	}

	*img = transform.Resize(*img, width, height, transform.Linear)
	return nil
}

// applyCrop crops the given image given the options
func applyCrop(img *image.Image, options options) error {
	width, height, x, y, err := options.Crop()
	if err != nil {
		return ErrEffectNotApplied
	}

	// calculate points
	x1 := x + width
	y1 := y + height

	*img = transform.Crop(*img, image.Rect(x, y, x1, y1))
	return nil
}
