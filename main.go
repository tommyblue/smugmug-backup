package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type smugMugConf struct {
	username    string
	destination string
	galleries   []gallery
	photos      []photo
}

type gallery struct{}
type photo struct{}

func main() {
	conf := parseArguments()

	// Check exising and writeability of destination folder
	conf.checkDestination()

	// Get user albums
	fmt.Printf("Getting albums for user %s...\n", conf.username)
	albums, err := conf.getAlbums()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d albums\n", len(*albums))

	// Iterate over all albums and:
	// - create folder
	// - iterate over all images
	//    - if existing, skip
	//    - if not, download
	for _, a := range *albums {
		createFolder(fmt.Sprintf("%s%s", conf.destination, a.URLPath))
		// a.Uris.AlbumImages.URI
	}
}

func (c *smugMugConf) checkDestination() {
	if err := checkFolderIsWritable(c.destination); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func parseArguments() *smugMugConf {
	conf := &smugMugConf{}

	flag.StringVar(&conf.username, "user", "", "SmugMug user to backup")
	flag.StringVar(&conf.destination, "destination", "", "Folder to save backup to")

	flag.Parse()

	if flag.NFlag() < 2 {
		fmt.Println("Missing arguments. Use --help for info")
		os.Exit(1)
	}

	return conf
}
