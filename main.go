package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var (
	// list of actions supported
	actionsFlag = map[string]func() error{
		"init":         SetupEnv,
		"delete-cache": deleteCache,
		"serve":        serve,
	}
)

func main() {

	var commName string
	if len(os.Args) < 2 {
		commName = "serve"
	} else {
		commName = os.Args[1]
	}

	comm, ok := actionsFlag[commName]
	if !ok {
		log.Fatalln(ErrInvalidFlag)
		return
	}

	if commName != "init" {
		loadEnv()
	}

	err := comm()
	if err != nil {
		log.Fatalln(err)
	}
}

// loadEnv loads the server.env file and puts the values
// in the env vars
func loadEnv() {
	err := godotenv.Load("server.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// serve starts the webapp
func serve() error {
	// setup default routes
	router := mux.NewRouter()
	router.HandleFunc("/", HandleNotAllowedMethod).Methods("GET")
	router.HandleFunc("/", HandleUpload).Methods("POST")
	router.HandleFunc("/{image}.{extension}", HandleDownload).Methods("GET")

	// serve
	serverPort := os.Getenv("SERVER_PORT")
	return http.ListenAndServe(":"+serverPort, router)
}
