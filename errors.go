package main

import (
	"errors"
	"net/http"
)

// Error codes returned by the download process
var (
	ErrInvalidExtension = errors.New("pix: invalid extension")
)

func fail(w http.ResponseWriter, err error, code int) {
	http.Error(w, err.Error(), code)
}
