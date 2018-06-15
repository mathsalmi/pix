package main

import (
	"errors"
	"net/http"
)

// Error codes returned by the download process
var (
	ErrInvalidExtension  = errors.New("pix: invalid extension")
	ErrOptionNotProvided = errors.New("pix: cannot parse values because this option was not provided")
	ErrInvalidDimensions = errors.New("pix: invalid image dimensions")
)

func fail(w http.ResponseWriter, err error, code int) {
	http.Error(w, err.Error(), code)
}

// HandleNotAllowedMethod writes Method Not Allowed to the response
// and sets the HTTP header status code to 405
func HandleNotAllowedMethod(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}
