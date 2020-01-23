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
	conf := parseArguments()

	// Check exising and writeability of destination folder
	conf.checkDestination()

	// Get user albums
	log.Infof("Getting albums for user %s...\n", conf.username)
	albums, err := conf.getAlbums()
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("Found %d albums\n", len(*albums))

	// Iterate over all albums and:
	// - create folder
	// - iterate over all images
	//    - if existing, skip
	//    - if not, download
	for _, a := range *albums {
		folder := fmt.Sprintf("%s%s", conf.destination, a.URLPath)
		createFolder(folder)
		images, err := conf.getAlbumImages(a.Uris.AlbumImages.URI)
		if err != nil {
			log.Fatal(err)
		}
		conf.saveImages(images, folder)
	}
}

func (c *smugMugConf) checkDestination() {
	if err := checkFolderIsWritable(c.destination); err != nil {
		log.Fatal(err)
	}
}

func parseArguments() *smugMugConf {
	conf := &smugMugConf{}

	flag.StringVar(&conf.username, "user", "", "SmugMug user to backup")
	flag.StringVar(&conf.destination, "destination", "", "Folder to save backup to")

	flag.Parse()

	if flag.NFlag() < 2 {
		log.Fatal("Missing arguments. Use --help for info")
	}

	return conf
}
