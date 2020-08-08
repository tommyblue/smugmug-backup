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
			s, err := smugmug.New(&smugmug.Conf{
				Username: "myUsername",
				Destination: "/path/to/backup/",
			})
			if err != nil {
				log.WithError(err).Fatal("Can't initialize the package")
			}

			if err := s.Run(); err != nil {
				log.Fatal(err)
			}
		}

	The package expects to find the SmugMug API credentials as
	environmental variables:

		API_KEY="<key>"
		API_SECRET="<secret>"
		USER_TOKEN="<Access Token>"
		USER_SECRET="<Token Secret>"
*/
package smugmug
