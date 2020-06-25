package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	log "github.com/sirupsen/logrus"

	"golang.org/x/sys/unix"
)

func createFolder(path string) {
	_, err := os.Stat(path)

	if err != nil {
		log.Infof("Creating folder %s\n", path)
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
			if !image.Processing { // Skip videos if under processing
				c.saveVideo(&image, folder)
			} else {
				log.Infof("Skipping video %s because under processing\n", image.Name())
			}
		} else {
			c.saveImage(&image, folder)
		}
	}
}

func (c *smugMugConf) saveImage(image *albumImage, folder string) {
	if image.Name() == "" {
		log.Warn("Unable to find valid file name, skipping..")
		return
	}
	dest := fmt.Sprintf("%s/%s", folder, image.Name())
	download(dest, image.ArchivedUri, image.ArchivedSize, c.ignorefetcherrors)
}

func (c *smugMugConf) saveVideo(image *albumImage, folder string) {
	if image.Name() == "" {
		log.Warn("Unable to find valid file name, skipping..")
		return
	}
	dest := fmt.Sprintf("%s/%s", folder, image.Name())

	var albumVideo albumVideo
	c.r.get(image.Uris.LargestVideo.Uri, &albumVideo)

	download(dest, albumVideo.Response.LargestVideo.Url, albumVideo.Response.LargestVideo.Size, c.ignorefetcherrors)
}

func sameFileSizes(path string, fileSize int64) bool {
	fi, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}
	return fi.Size() == fileSize
}

func download(dest, downloadURL string, fileSize int64, ignorefetcherrors bool) {
	if _, err := os.Stat(dest); err == nil {
		if sameFileSizes(dest, fileSize) {
			log.Debugf("File exists with same size: %s\n", downloadURL)
			return
		}
	}
	log.Infof("Getting %s\n", downloadURL)

	response, err := makeAPICall(downloadURL)
	if err != nil {
		if ignorefetcherrors {
			log.Errorf("%s: download failed with: %s", downloadURL, err)
			return
		}
		log.Fatalf("%s: download failed with: %s", downloadURL, err)
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
	log.Infof("Saved %s\n", dest)
}
