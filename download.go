package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/gorilla/mux"
)

// HandleDownload writes the requested image to the response.
//
// It may optionally apply effects and transformations on images, as
// requested in the URL.
//
// In case a file does not exist or any I/O error occurs, it will write
// an HTTP error response, like 404 (not found) or 500 (internal server error).
func HandleDownload(w http.ResponseWriter, r *http.Request) {
	var (
		uploadDir = os.Getenv("UPLOAD_DIR")

		vars      = mux.Vars(r)
		filename  = vars["image"]
		extension = vars["extension"]
	)

	err := validateExtension(extension)
	if err != nil {
		fail(w, err, http.StatusInternalServerError)
		return
	}

	filepaths, err := filepath.Glob(fmt.Sprintf("%s/%s.*", uploadDir, filename))
	if err != nil || filepaths == nil {
		http.NotFound(w, r)
		return
	}

	filepath := filepaths[0]

	originalImage, err := imgio.Open(filepath)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	img := &Image{originalImage}

	options := parseOptions(r, extension, img)
	newpath := createNewPath(filename, extension, options)

	if !isCached(newpath) {
		ApplyTransformations(img, options)

		if err := imgio.Save(newpath, img, options.Encoder()); err != nil {
			fail(w, err, http.StatusInternalServerError)
			return
		}
	}

	// deliver
	data, err := ioutil.ReadFile(newpath)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Header().Set("Content-Type", http.DetectContentType(data))

	w.Write(data)
}

// validateExtension returns ErrInvalidExtension if the given extension is
// invalid or nil if it is valid.
func validateExtension(extension string) error {
	exts := []string{"jpeg", "jpg", "png", "gif", "bmp"}

	for _, e := range exts {
		if e == extension {
			return nil
		}
	}

	return ErrInvalidExtension
}

// createNewPath returns the filename for an image given its name, extension and options
func createNewPath(filename, extension string, options *Options) string {
	uploadDir := os.Getenv("UPLOAD_DIR")
	return fmt.Sprintf("%s/%s-%s.%s", uploadDir, filename, options.Hash(), extension)
}

// isCached returns true if image already exists or false if it does not
func isCached(filepath string) bool {
	_, err := os.Stat(filepath)
	return !os.IsNotExist(err)
}
