package main

import (
	"io"
	"net/http"
)

func handleUpload(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "This handles upload")
}
