package main

import (
	"flag"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

type smugMugConf struct {
	username    string
	destination string
	req         requestsHandler
}

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)

	flag, isPresent := os.LookupEnv("DEBUG")
	if isPresent && flag == "1" {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)

	}
}

func main() {
	conf := &smugMugConf{
		req: &smugmugHandler{},
	}
	conf.parseArguments()

	// Check exising and writeability of destination folder
	if err := checkFolderIsWritable(conf.destination); err != nil {
		log.Fatal("Can't write in the destination folder")
	}

	// Get user albums
	log.Infof("Getting albums for user %s...\n", conf.username)
	uri := conf.getUserAlbumsURI()
	albums, err := conf.getAlbums(uri)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("Found %d albums\n", len(albums))

	// Iterate over all albums and:
	// - create folder
	// - iterate over all images
	//    - if existing, skip
	//    - if not, download
	var errors int
	for _, album := range albums {
		folder := fmt.Sprintf("%s%s", conf.destination, album.URLPath)

		if err := createFolder(folder); err != nil {
			log.Errorf("cannot create the destination folder %s. Error: %v", folder, err)
			errors++
			continue
		}

		log.Debugf("[ALBUM IMAGES] %s", album.Uris.AlbumImages.URI)
		images, err := conf.getAlbumImages(album.Uris.AlbumImages.URI, album.URLPath)
		if err != nil {
			log.Errorf("Cannot get album images for %s. Error: %v", album.Uris.AlbumImages.URI, err)
			errors++
			continue
		}
		log.Debugf("Got album images for %s", album.Uris.AlbumImages.URI)
		log.Debugf("%+v", images)
		conf.saveImages(images, folder)
	}

	if errors > 0 {
		log.Fatalf("Completed with %d errors, check logs", errors)
	}

	log.Info("Backup completed")
}

func (c *smugMugConf) parseArguments() {
	flag.StringVar(&c.username, "user", "", "SmugMug user to backup")
	flag.StringVar(&c.destination, "destination", "", "Folder to save backup to")

	flag.Parse()

	if flag.NFlag() < 2 {
		log.Fatal("Missing arguments. Use --help for info")
	}
}
