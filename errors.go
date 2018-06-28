package main

import (
	"errors"
	"net/http"
)

// General error codes
var (
	ErrNotImplemented = errors.New("pix: not implemented")
	ErrInternalServer = errors.New("pix: server is improperly set up")
)

// Error codes used in CLI flag parsing
var (
	ErrInvalidFlag         = errors.New("pix: invalid flag provided")
	ErrCacheNoFilesDeleted = errors.New("pix: no files deleted")
	ErrSetupEnvFile        = errors.New("pix: error creating the `server.env` file")
)

// Error codes returned by the download process
var (
	ErrInvalidExtension = errors.New("pix: invalid extension")
)

// Error codes for options parsing
var (
	ErrOptionNotProvided   = errors.New("pix: option not provided")
	ErrInvalidOptionValues = errors.New("pix: cannot parse option because the provided values are invalid")
	ErrInvalidDimensions   = errors.New("pix: invalid image dimensions")
)

// Error codes for effects/transformation process
var (
	ErrTransformationNotApplied = errors.New("pix: effect not applied")
)

// Error codes used in the upload process
var (
	ErrInvalidUpload   = errors.New("pix: invalid upload")
	ErrFileTooBig      = errors.New("pix: upload: file too big")
	ErrInvalidFileType = errors.New("pix: upload: invalid file type")
)

// fail writes the string representation of an error to the response
func fail(w http.ResponseWriter, err error, code int) {
	http.Error(w, err.Error(), code)
}

// HandleNotAllowedMethod writes Method Not Allowed to the response
// and sets the HTTP header status code to 405
func HandleNotAllowedMethod(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}
