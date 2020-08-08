package main

import (
	"errors"
	"flag"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/tommyblue/smugmug-backup"
)

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
	cfg, err := parseArguments()
	if err != nil {
		log.WithError(err).Fatal("Can't parse command line arguments")
	}

	s, err := smugmug.New(cfg)
	if err != nil {
		log.WithError(err).Fatal("Can't initialize the package")
	}

	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}

func parseArguments() (*smugmug.Conf, error) {
	cfg := &smugmug.Conf{}
	flag.StringVar(&cfg.Username, "user", "", "SmugMug user to backup")
	flag.StringVar(&cfg.Destination, "destination", "", "Folder to save backup to")

	flag.Parse()

	if flag.NFlag() < 2 {
		return nil, errors.New("Missing arguments. Use --help for info")
	}

	if cfg.Username == "" {
		return nil, errors.New("-user is a required argument")
	}

	if cfg.Destination == "" {
		return nil, errors.New("-destination is a required argument")
	}

	return cfg, nil
}
