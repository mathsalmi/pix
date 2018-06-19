package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"image"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/anthonynsimon/bild/imgio"
)

// parseOptions returns the transformation options extracted
// from the request
func parseOptions(r *http.Request, extension string, img *image.Image) options {

	options := make(options)

	values := r.URL.Query()
	for key := range values {
		options[key] = values.Get(key)
	}

	// set standard options
	options["extension"] = extension

	width := (*img).Bounds().Max.X
	height := (*img).Bounds().Max.Y
	options["original_width"] = strconv.Itoa(width)
	options["original_height"] = strconv.Itoa(height)

	return options
}

// options holds the transformations settings
// to be applied in the requested image
type options map[string]string

// Extension returns the desired extension
func (o options) Extension() string {
	return o["extension"]
}

// Hash returns a unique string MD5 encoded that represents the
// transformation options provided.
func (o options) Hash() string {
	b := new(bytes.Buffer)

	s := make([]string, 0, len(o))
	for key := range o {
		s = append(s, key)
	}

	sort.Strings(s)

	for _, key := range s {
		value := o[key]
		b.WriteString(key + value)
	}

	return fmt.Sprintf("%x", md5.Sum(b.Bytes()))
}

// Quality returns the picture quality for JPEG images.
//
// If the provided value is an invalid number, less than 1,
// or greater 100, 80 is returned instead.
func (o options) Quality() int {
	value, ok := o["quality"]
	if !ok {
		return 80
	}

	number, err := strconv.Atoi(value)
	if err != nil || number < 1 || number > 100 {
		return 80
	}

	return number
}

// NumColors returns the maximum number of colors used in GIF images.
//
// It ranges from 1 to 256.
func (o options) NumColors() int {
	value, ok := o["numcolors"]
	if !ok {
		return 256
	}

	number, err := strconv.Atoi(value)
	if err != nil || number < 1 || number > 256 {
		return 256
	}

	return number
}

// Encoder returns the image encoder accordingly to the desired
// image extension
func (o options) Encoder() imgio.Encoder {
	switch o.Extension() {
	case "jpg", "jpeg":
		return imgio.JPEGEncoder(o.Quality())
	case "png":
		return imgio.PNGEncoder()
	case "bmp":
		return imgio.BMPEncoder()
	case "gif":
		return GIFEncoder(o.NumColors())
	}
	return nil
}

// Resize calculates the new values for resizing the image.
func (o options) Resize() (int, int, error) {
	oWidth, hasWidth := o["width"]
	oHeight, hasHeight := o["height"]

	ooWidth, _ := o["original_width"]
	ooHeight, _ := o["original_height"]

	if !hasWidth && !hasHeight {
		return 0, 0, ErrOptionNotProvided
	}

	width, height := 0, 0
	originalWidth, _ := strconv.Atoi(ooWidth)
	originalHeight, _ := strconv.Atoi(ooHeight)

	// calculate values
	if hasWidth && hasHeight {
		width, _ = strconv.Atoi(oWidth)
		height, _ = strconv.Atoi(oHeight)
	} else if hasWidth && !hasHeight {
		width, _ = strconv.Atoi(oWidth)
		height = (width * originalHeight) / originalWidth
	} else if !hasWidth && hasHeight {
		height, _ = strconv.Atoi(oHeight)
		width = (height * originalWidth) / originalHeight
	}

	// check boundaries
	if width < 0 || height < 0 {
		return 0, 0, ErrInvalidDimensions
	}

	// TODO(salmi): put a max file width/height check here?

	return width, height, nil
}

// Crop checks the values for cropping the image.
//
// It has to be applied to the original image, so the execution
// order of transformation functions matters in this case.
func (o options) Crop() (width int, height int, x int, y int, err error) {
	s, ok := o["crop"]
	if !ok {
		err = ErrOptionNotProvided
		return
	}

	values, e := parseCropString(s)
	if e != nil {
		err = ErrInvalidOptionValues
		return
	}

	width, hasWidth := values["w"]
	height, hasHeight := values["h"]
	x, hasX := values["x"]
	y, hasY := values["y"]

	if !hasWidth || !hasHeight || !hasX || !hasY {
		err = ErrInvalidOptionValues
		return
	}

	originalWidth := o["original_width"]
	originalHeight := o["original_height"]

	// check boundaries
	if w, _ := strconv.Atoi(originalWidth); x > w || x+width > w {
		err = ErrInvalidOptionValues
		return
	}

	if h, _ := strconv.Atoi(originalHeight); y > h || y+height > h {
		err = ErrInvalidOptionValues
		return
	}

	return
}

// SmartCrop gets the options to apply the smart crop algorithm
//
// It has to be applied to the original image, so the execution
// order of transformation functions matters in this case.
func (o options) SmartCrop() (width, height int, err error) {
	s, ok := o["smartcrop"]
	if !ok {
		err = ErrOptionNotProvided
		return
	}

	values, e := parseCropString(s)
	if e != nil {
		err = ErrInvalidOptionValues
		return
	}

	width, hasWidth := values["w"]
	height, hasHeight := values["h"]

	if !hasWidth || !hasHeight {
		err = ErrInvalidOptionValues
		return
	}

	originalWidth := o["original_width"]
	originalHeight := o["original_height"]

	// check boundaries
	if w, _ := strconv.Atoi(originalWidth); width > w {
		err = ErrInvalidOptionValues
		return
	}

	if h, _ := strconv.Atoi(originalHeight); height > h {
		err = ErrInvalidOptionValues
		return
	}

	return
}

// FlipH tells whether or not to apply the horizontal
// flip transformation
func (o options) FlipH() error {
	_, ok := o["flipH"]
	if !ok {
		return ErrOptionNotProvided
	}

	return nil
}

// FlipV tells whether or not to apply the vertical
// flip transformation
func (o options) FlipV() error {
	_, ok := o["flipV"]
	if !ok {
		return ErrOptionNotProvided
	}

	return nil
}

// parseCropString parses the crop string and returns
// a map contain the values.
//
// The string looks like: `w:300|h:300|x:20|y:30` and
// it may contain less or more values depending on the
// context it is being applied
func parseCropString(s string) (values map[string]int, err error) {
	values = make(map[string]int)

	pairs := strings.Split(s, "|")
	for _, pair := range pairs {
		v := strings.Split(pair, ":")
		if len(v) != 2 {
			err = ErrInvalidOptionValues
			return
		}

		key := v[0]
		value, e := strconv.Atoi(v[1])
		if e != nil {
			err = ErrInvalidOptionValues
			return
		}

		// check boundaries - fast fail
		if value < 0 {
			err = ErrInvalidOptionValues
			return
		}

		values[key] = value
	}

	return
}
