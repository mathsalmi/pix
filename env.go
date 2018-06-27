package main

import (
	"github.com/joho/godotenv"
)

// DefaultServerOptions provides the values for creating
// a default `server.env` file
var DefaultServerOptions = map[string]string{
	"SERVER_PORT": "8000",
	"UPLOAD_DIR":  "~/pix/upload",
}

// SetupEnv creates all files and directories needed to run Pix
func SetupEnv() error {
	err := createServerConfFile()
	if err != nil {
		return err
	}
	return nil
}

// createServerConfFile creates an example `server.env` file with
// the default values provided in DefaultServerOptions
func createServerConfFile() error {
	err := godotenv.Write(DefaultServerOptions, "./server.env")
	if err != nil {
		return ErrSetupEnvFile
	}
	return nil
}
