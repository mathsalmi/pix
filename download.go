package main

import (
	"io"
	"net/http"
)

func handleDownload(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Handle downloads")
}
