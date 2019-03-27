package main

import (
	"flag"
	"fmt"
	"os"
)

type smugMugConf struct {
	apiKey      string
	apiSecret   string
	username    string
	destination string
	galleries   []gallery
	photos      []photo
}

type gallery struct{}
type photo struct{}

func main() {
	conf := parseArguments()

	conf.checkDestination()
	// Check exising and writeability of destination folder
	// Get user root node
	// Iterate over all nodes and collect galleries
	// For each gallery, collect photos info
	// Verify if already existing, download if not

}

func (c *smugMugConf) checkDestination() {
	if err := checkFolderIsWritable(c.destination); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func parseArguments() *smugMugConf {
	conf := &smugMugConf{}

	flag.StringVar(&conf.apiKey, "key", "", "API Key")
	flag.StringVar(&conf.apiSecret, "secret", "", "API Secret")
	flag.StringVar(&conf.username, "user", "", "SmugMug user to backup")
	flag.StringVar(&conf.destination, "destination", "", "Folder to save backup to")

	flag.Parse()

	if flag.NFlag() < 4 {
		fmt.Println("Missing arguments. Use --help for info")
		os.Exit(1)
	}

	return conf
}
