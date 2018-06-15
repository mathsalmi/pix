package main

import (
	"io"
	"net/http"
)

// HandleUpload handles image uploads requests
func HandleUpload(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "This handles upload")
}
