package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
)

func handleDownload(w http.ResponseWriter, r *http.Request) {
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

	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Header().Set("Content-Type", http.DetectContentType(data))

	w.Write(data)
}

func validateExtension(extension string) error {
	exts := []string{"jpeg", "jpg", "png", "gif"}

	for _, e := range exts {
		if e == extension {
			return nil
		}
	}

	return ErrInvalidExtension
}
