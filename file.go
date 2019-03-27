package main

import (
	"errors"
	"os"
	"path"

	"golang.org/x/sys/unix"
)

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
