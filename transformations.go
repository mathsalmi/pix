package main

import (
	"image"

	"github.com/anthonynsimon/bild/transform"
	"github.com/muesli/smartcrop"
	"github.com/muesli/smartcrop/nfnt"
)

// ApplyTransformations applies transformations to the given image
func ApplyTransformations(img *image.Image, options *Options) {
	applyCrop(img, options)
	applySmartCrop(img, options)
	applyResize(img, options)
	applyFlipH(img, options)
	applyFlipV(img, options)
}

// applyResize resizes the given image given the options
func applyResize(img *image.Image, options *Options) error {
	width, height, err := options.Resize()
	if err != nil {
		return ErrTransformationNotApplied
	}

	*img = transform.Resize(*img, width, height, transform.Linear)
	return nil
}

// applyCrop crops the given image given the options
func applyCrop(img *image.Image, options *Options) error {
	width, height, x, y, err := options.Crop()
	if err != nil {
		return ErrTransformationNotApplied
	}

	// calculate points
	x1 := x + width
	y1 := y + height

	*img = transform.Crop(*img, image.Rect(x, y, x1, y1))
	return nil
}

// applySmartCrop crops the image using the Smart Crop algorithm
// applying the provided options
func applySmartCrop(img *image.Image, options *Options) error {
	width, height, err := options.SmartCrop()
	if err != nil {
		return ErrTransformationNotApplied
	}

	analyzer := smartcrop.NewAnalyzer(nfnt.NewDefaultResizer())
	rect, err := analyzer.FindBestCrop(*img, width, height)
	if err != nil {
		return ErrTransformationNotApplied
	}

	*img = transform.Crop(*img, rect)
	return nil
}

// applyFlipH flips image horizontally
func applyFlipH(img *image.Image, options *Options) error {
	err := options.FlipH()
	if err != nil {
		return ErrTransformationNotApplied
	}

	*img = transform.FlipH(*img)
	return nil
}

// applyFlipV flips image vertically
func applyFlipV(img *image.Image, options *Options) error {
	err := options.FlipV()
	if err != nil {
		return ErrTransformationNotApplied
	}

	*img = transform.FlipV(*img)
	return nil
}
