package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// isCached returns true if image already exists or false if it does not
func isCached(filepath string) bool {
	_, err := os.Stat(filepath)
	return !os.IsNotExist(err)
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
