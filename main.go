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
	// setup default routes
	router := mux.NewRouter()
	router.HandleFunc("/", HandleNotAllowedMethod).Methods("GET")
	router.HandleFunc("/", HandleUpload).Methods("POST")
	router.HandleFunc("/{image}.{extension}", HandleDownload).Methods("GET")

	// serve
	serverPort := os.Getenv("SERVER_PORT")
	log.Fatal(http.ListenAndServe(":"+serverPort, router))
}
