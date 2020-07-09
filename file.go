package main

import (
	"errors"
	"fmt"
	"os"
	"path"

	log "github.com/sirupsen/logrus"

	"golang.org/x/sys/unix"
)

func createFolder(path string) error {
	_, err := os.Stat(path)

	// Folder exists
	if err == nil {
		return nil
	}

	log.Infof("Creating folder %s\n", path)
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return fmt.Errorf("Cannot create folder: %v", err)
	}

	return nil
}

func checkFolderIsWritable(folderPath string) error {
	if !path.IsAbs(folderPath) {
		return errors.New("Destination path must be an absolute path")
	}

	info, err := os.Stat(folderPath)
	if err != nil {
		return errors.New("Destination path doesn't exist")
	}

	if !info.IsDir() {
		return errors.New("Destination path isn't a directory")
	}

	if unix.Access(folderPath, unix.W_OK) != nil {
		return errors.New("Destination path in not writeable")
	}

	return nil
}

func sameFileSizes(path string, fileSize int64) bool {
	fi, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}
	return fi.Size() == fileSize
}
