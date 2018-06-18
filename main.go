package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var (
	// list of actions supported
	actionsFlag = map[string]func() error{
		"init":         setupEnv,
		"delete-cache": deleteCache,
		"serve":        serve,
	}
)

func init() {
	// Load env file
	err := godotenv.Load("server.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {

	if len(os.Args) < 2 {
		serve()
		return
	}

	comm, ok := actionsFlag[os.Args[1]]
	if !ok {
		log.Fatalln(ErrInvalidFlag.Error())
		return
	}

	err := comm()
	if err != nil {
		log.Fatalln(err.Error())
	}
}

// deleteCache deleted all cached images
func deleteCache() error {
	uploadDir := os.Getenv("UPLOAD_DIR")

	filepaths, err := filepath.Glob(fmt.Sprintf("%s/*-*.*", uploadDir))
	if err != nil {
		return ErrCacheNoFilesDeleted
	}

	for _, file := range filepaths {
		if err := os.Remove(file); err != nil {
			return err
		}
	}

	return nil
}

// setupEnv creates all files and directories needed to run Pix
func setupEnv() error {
	return ErrNotImplemented
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
