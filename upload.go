package main

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"strconv"

	"github.com/segmentio/ksuid"
)

// HandleUpload handles image uploads requests
func HandleUpload(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("image")
	if err != nil {
		fail(w, ErrInvalidUpload, http.StatusBadRequest)
		return
	}
	defer file.Close()

	umax := os.Getenv("MAX_UPLOAD_SIZE")
	maxUploadSizeAllowed, err := strconv.Atoi(umax)
	if err != nil {
		fail(w, ErrInternalServer, http.StatusInternalServerError)
		return
	}

	if handler.Size > int64(maxUploadSizeAllowed) {
		fail(w, ErrFileTooBig, http.StatusRequestEntityTooLarge)
		return
	}

	header := handler.Header
	mimeType := header.Get("Content-Type")
	if err := validateFileType(mimeType); err != nil {
		fail(w, err, http.StatusBadRequest)
		return
	}

	var newextension string
	if extensions, err := mime.ExtensionsByType(mimeType); err != nil {
		fail(w, err, http.StatusInternalServerError)
		return
	} else {
		newextension = extensions[0]
	}

	newname := ksuid.New().String() + newextension
	newfilepath := fmt.Sprintf("%s/%s", os.Getenv("UPLOAD_DIR"), newname)

	newfile, err := os.OpenFile(newfilepath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		fail(w, err, http.StatusInternalServerError)
		return
	}
	defer newfile.Close()

	if _, err := io.Copy(newfile, file); err != nil {
		fail(w, err, http.StatusInsufficientStorage)
		return
	}

	io.WriteString(w, newname)
}

// validateFileType returns ErrInvalidFileType if given mime type is invalid
func validateFileType(mime string) error {
	types := []string{"image/jpeg", "image/jpg", "image/gif", "image/png", "image/bmp"}
	for _, t := range types {
		if t == mime {
			return nil
		}
	}

	return ErrInvalidFileType
}
