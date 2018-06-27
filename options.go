package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/anthonynsimon/bild/imgio"
)

// parseOptions returns the transformation options extracted
// from the request
func parseOptions(r *http.Request, extension string, img *Image) *Options {

	options := NewOptions()

	values := r.URL.Query()
	for key := range values {
		v, err := strconv.Atoi(values.Get(key))
		if err != nil {
			options.SetString(key, values.Get(key))
		} else {
			options.SetInt(key, v)
		}
	}

	// set standard options
	options.SetString("extension", extension)
	options.img = img

	return options
}

// Options holds the transformations settings
// to be applied in the requested image
type Options struct {
	values map[string]interface{}
	img    *Image
}

// NewOptions returns a new options with values
func NewOptions() *Options {
	return &Options{
		values: make(map[string]interface{}),
	}
}

// SetInt inserts an int value given its key in options
func (o *Options) SetInt(key string, value int) {
	o.values[key] = value
}

// SetString inserts a string value given its key in options
func (o *Options) SetString(key, value string) {
	o.values[key] = value
}

// Int returns an int value.Int
//
// ok returns false if the key does not exist or
// if the value cannot be asserted into int type
func (o Options) Int(key string) (value int, ok bool) {
	v, ok := o.values[key]
	if !ok {
		return 0, false
	}

	value, ok = v.(int)
	return value, ok
}

// String returns a string value
//
// ok returns false if the key does not exist or
// if the value cannot be asserted into string type
func (o Options) String(key string) (value string, ok bool) {
	v, ok := o.values[key]
	if !ok {
		return "", false
	}

	value, ok = v.(string)
	return value, ok
}

// Image returns the image under modification
func (o Options) Image() *Image {
	return o.img
}

// Extension returns the desired extension
func (o Options) Extension() string {
	v, _ := o.String("extension")
	return v
}

// Hash returns a unique string MD5 encoded that represents the
// transformation options provided.
func (o Options) Hash() string {
	b := new(bytes.Buffer)

	s := make([]string, 0, len(o.values))
	for key := range o.values {
		s = append(s, key)
	}

	sort.Strings(s)

	for _, key := range s {
		value, _ := o.String(key)
		b.WriteString(key + value)
	}

	return fmt.Sprintf("%x", md5.Sum(b.Bytes()))
}

// Quality returns the picture quality for JPEG images.
//
// If the provided value is an invalid number, less than 1,
// or greater 100, 80 is returned instead.
func (o Options) Quality() int {
	number, ok := o.Int("quality")
	if !ok || number < 1 || number > 100 {
		return 80
	}

	return number
}

// NumColors returns the maximum number of colors used in GIF images.
//
// It ranges from 1 to 256.
func (o Options) NumColors() int {
	number, ok := o.Int("numcolors")
	if !ok || number < 1 || number > 256 {
		return 256
	}

	return number
}

// Encoder returns the image encoder accordingly to the desired
// image extension
func (o Options) Encoder() imgio.Encoder {
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
func (o Options) Resize() (width, height int, err error) {
	width, hasWidth := o.Int("width")
	height, hasHeight := o.Int("height")

	originalWidth := o.Image().Width()
	originalHeight := o.Image().Height()

	if !hasWidth && !hasHeight {
		return 0, 0, ErrOptionNotProvided
	}

	// calculate values
	if hasWidth && !hasHeight {
		height = (width * originalHeight) / originalWidth
	} else if !hasWidth && hasHeight {
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
func (o Options) Crop() (width, height, x, y int, err error) {
	s, ok := o.String("crop")
	if !ok {
		return 0, 0, 0, 0, ErrOptionNotProvided
	}

	values, e := parseCropString(s)
	if e != nil {
		return 0, 0, 0, 0, ErrInvalidOptionValues
	}

	width, hasWidth := values["w"]
	height, hasHeight := values["h"]
	x, hasX := values["x"]
	y, hasY := values["y"]

	if !hasWidth || !hasHeight || !hasX || !hasY {
		return 0, 0, 0, 0, ErrInvalidOptionValues
	}

	originalWidth := o.Image().Width()
	originalHeight := o.Image().Height()

	// check boundaries
	if x > originalWidth || x+width > originalWidth {
		return 0, 0, 0, 0, ErrInvalidOptionValues
	}

	if y > originalHeight || y+height > originalHeight {
		return 0, 0, 0, 0, ErrInvalidOptionValues
	}

	return width, height, x, y, nil
}

// SmartCrop gets the options to apply the smart crop algorithm
//
// It has to be applied to the original image, so the execution
// order of transformation functions matters in this case.
func (o Options) SmartCrop() (width, height int, err error) {
	s, ok := o.String("smartcrop")
	if !ok {
		return 0, 0, ErrOptionNotProvided
	}

	values, e := parseCropString(s)
	if e != nil {
		return 0, 0, ErrInvalidOptionValues
	}

	width, hasWidth := values["w"]
	height, hasHeight := values["h"]

	if !hasWidth || !hasHeight {
		return 0, 0, ErrInvalidOptionValues
	}

	originalWidth := o.Image().Width()
	originalHeight := o.Image().Height()

	// check boundaries
	if width > originalWidth {
		return 0, 0, ErrInvalidOptionValues
	}

	if height > originalHeight {
		return 0, 0, ErrInvalidOptionValues
	}

	return width, height, nil
}

// FlipH tells whether or not to apply the horizontal
// flip transformation
func (o Options) FlipH() error {
	_, ok := o.String("flipH")
	if !ok {
		return ErrOptionNotProvided
	}

	return nil
}

// FlipV tells whether or not to apply the vertical
// flip transformation
func (o Options) FlipV() error {
	_, ok := o.String("flipV")
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
			return nil, ErrInvalidOptionValues
		}

		key := v[0]
		value, e := strconv.Atoi(v[1])
		if e != nil {
			return nil, ErrInvalidOptionValues
		}

		// check boundaries - fast fail
		if value < 0 {
			return nil, ErrInvalidOptionValues
		}

		values[key] = value
	}

	return values, nil
}
