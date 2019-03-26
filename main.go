package main

import (
	"flag"
	"fmt"
)

type smugMugConf struct {
	apiKey    string
	apiSecret string
	username  string
}

func main() {
	conf := parseArguments()
}

func parseArguments() *smugMugConf {
	conf := &smugMugConf{}

	flag.StringVar(&conf.apiKey, "key", "", "API Key")
	flag.StringVar(&conf.apiSecret, "secret", "", "API Secret")
	flag.StringVar(&conf.username, "user", "", "SmugMug user to backup")

	flag.Parse()

	if flag.NFlag() < 3 {
		fmt.Println("Missing arguments. Use --help for info")
	}

	return conf
}
