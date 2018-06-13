package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// setup routes
	router := mux.NewRouter()
	router.HandleFunc("/", handleIndex).Methods("GET")
	router.HandleFunc("/", handleUpload).Methods("POST")
	router.HandleFunc("/{image}.{extension}", handleDownload).Methods("GET")

	// serve
	log.Fatal(http.ListenAndServe(":8000", router))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}
