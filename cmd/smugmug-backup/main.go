package main

import (
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
	cfg, err := smugmug.ReadConf()
	if err != nil {
		log.WithError(err).Fatal("Configuration error")
	}

	wrk, err := smugmug.New(cfg)
	if err != nil {
		log.WithError(err).Fatal("Can't initialize the package")
	}

	if err := wrk.Run(); err != nil {
		log.Fatal(err)
	}
}
