/*
	Package smugmug implements the logic to perform a full
	backup of a SmugMug account (images and videos).

	The program loops over the images and videos of the user's
	albums and saves them in the destination folder, replicating
	the SmugMug paths.

	You can run the app multiple times, all existing files will
	be skipped if their sizes match.

	Creating a binary with this package is as simple as:

		package main

		import (
			log "github.com/sirupsen/logrus"
			"github.com/tommyblue/smugmug-backup"
		)

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

	The app reads its configuration from ./config.toml or $HOME/.smgmg/config.toml.

	The supported configuration keys/values are the following:

		[authentication]
		username = "<SmugMug username>"
		api_key = "<API Key>"
		api_secret = "<API Secret>"
		user_token = "<User Token>"
		user_secret = "<User Secret>"

		[store]
		destination = "<Backup destination folder>"

	All values can be overridden by environment variables, that have the following names:

		SMGMG_BK_USERNAME = "<SmugMug username>"
		SMGMG_BK_API_KEY = "<API Key>"
		SMGMG_BK_API_SECRET = "<API Secret>"
		SMGMG_BK_USER_TOKEN = "<User Token>"
		SMGMG_BK_USER_SECRET = "<User Secret>"
		SMGMG_BK_DESTINATION = "<Backup destination folder>"

	All configuration values are required. They can be omitted in the configuration file
	as long as they are overridden by environment values.
*/
package smugmug
