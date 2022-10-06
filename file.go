package smugmug

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func createFolder(path string) error {
	_, err := os.Stat(path)

	// Folder exists
	if err == nil {
		return nil
	}

	log.Infof("Creating folder %s\n", path)
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return fmt.Errorf("cannot create folder: %v", err)
	}

	return nil
}

func checkDestFolder(folderPath string) error {
	if !filepath.IsAbs(folderPath) {
		return errors.New("destination path must be an absolute path")
	}

	info, err := os.Stat(folderPath)
	if err != nil {
		return errors.New("destination path doesn't exist")
	}

	if !info.IsDir() {
		return errors.New("destination path isn't a directory")
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
