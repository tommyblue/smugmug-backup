package smugmug

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// Conf is the configuration of the smugmug worker
type Conf struct {
	Username    string
	Destination string
}

// Worker actually implements the backup logic
type Worker struct {
	req requestsHandler
	cfg *Conf
}

// New return a SmugMug backup configuration. It returns an error if it fails parsing
// the command line arguments
func New(cfg *Conf) (*Worker, error) {
	if cfg.Username == "" {
		return nil, errors.New("Username can't be empty")
	}

	if cfg.Destination == "" {
		return nil, errors.New("Destination can't be empty")
	}

	// Check exising and writeability of destination folder
	if err := checkDestFolder(cfg.Destination); err != nil {
		return nil, fmt.Errorf("Can't find in the destination folder: %v", err)
	}

	return &Worker{
		cfg: cfg,
		req: &handler{},
	}, nil
}

// Run runs the backup of the whole account
func (w *Worker) Run() error {
	// Get user albums
	log.Infof("Getting albums for user %s...\n", w.cfg.Username)
	albums, err := w.userAlbums()
	if err != nil {
		return fmt.Errorf("Error getting user albums: %v", err)
	}

	log.Infof("Found %d albums\n", len(albums))

	// Iterate over all albums and:
	// - create folder
	// - iterate over all images/videos
	//    - if existing, skip
	//    - if not, download
	var errors int
	for _, album := range albums {
		folder := fmt.Sprintf("%s%s", w.cfg.Destination, album.URLPath)

		if err := createFolder(folder); err != nil {
			log.WithError(err).Errorf("cannot create the destination folder %s", folder)
			errors++
			continue
		}

		log.Debugf("[ALBUM IMAGES] %s", album.Uris.AlbumImages.URI)
		images, err := w.albumImages(album.Uris.AlbumImages.URI, album.URLPath)
		if err != nil {
			log.WithError(err).Errorf("Cannot get album images for %s", album.Uris.AlbumImages.URI)
			errors++
			continue
		}
		log.Debugf("Got album images for %s", album.Uris.AlbumImages.URI)
		log.Debugf("%+v", images)
		w.saveImages(images, folder)
	}

	if errors > 0 {
		return fmt.Errorf("Completed with %d errors, check logs", errors)
	}

	log.Info("Backup completed")
	return nil
}
