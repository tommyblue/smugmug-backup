package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"golang.org/x/sys/unix"
)

func createFolder(path string) {
	_, err := os.Stat(path)

	if err != nil {
		log.Printf("Creating folder %s", path)
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			panic("Cannot create folder")
		}
	}
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

func (c *smugMugConf) saveImages(images *[]albumImage, folder string) {
	for _, image := range *images {
		if image.IsVideo {
			c.saveVideo(&image, folder)
		} else {
			c.saveImage(&image, folder)
		}
	}
}

func (c *smugMugConf) saveImage(image *albumImage, folder string) {
	dest := fmt.Sprintf("%s/%s", folder, image.FileName)
	download(dest, image.ArchivedUri, image.ArchivedSize)
}

func (c *smugMugConf) saveVideo(image *albumImage, folder string) {
	dest := fmt.Sprintf("%s/%s", folder, image.FileName)

	var albumVideo albumVideo
	c.get(image.Uris.LargestVideo.Uri, &albumVideo)

	download(dest, albumVideo.Response.LargestVideo.Url, albumVideo.Response.LargestVideo.Size)
}

func sameFileSizes(path string, fileSize int64) bool {
	fi, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}
	return fi.Size() == fileSize
}

func download(dest, downloadURL string, fileSize int64) {

	if _, err := os.Stat(dest); err == nil {
		if sameFileSizes(dest, fileSize) {
			debugMsg(fmt.Sprintf("File exists with same size: %s\n", downloadURL))
			return
		}
	}
	logMsg(fmt.Sprintf("Getting %s\n", downloadURL))

	response, err := makeAPICall(downloadURL)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	file, err := os.Create(dest)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}
	logMsg(fmt.Sprintf("Saved %s\n\n", dest))
}
