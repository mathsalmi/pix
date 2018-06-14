package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	// setup routes
	router := mux.NewRouter()
	router.HandleFunc("/", handleIndex).Methods("GET")
	router.HandleFunc("/", handleUpload).Methods("POST")
	router.HandleFunc("/{image}.{extension}", handleDownload).Methods("GET")

	// serve
	serverPort := os.Getenv("SERVER_PORT")
	log.Fatal(http.ListenAndServe(":"+serverPort, router))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}
