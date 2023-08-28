package smugmug

import (
	"errors"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

type requestsHandler interface {
	get(string, interface{}) error
}

// userAlbums returns the list of albums belonging to the suer
func (w *Worker) userAlbums() ([]album, error) {
	uri := w.userAlbumsURI()
	return w.albums(uri)
}

// currentUser returns the nickname of the authenticated user
func (w *Worker) currentUser() (string, error) {
	var u currentUser
	if err := w.req.get("/api/v2!authuser", &u); err != nil {
		return "", err
	}
	return u.Response.User.NickName, nil
}

// userAlbumsURI returns the URI of the first page of the user albums. It's intended to be used
// as argument for a call to albums()
func (w *Worker) userAlbumsURI() string {
	var u user
	path := fmt.Sprintf("/api/v2/user/%s", w.cfg.username)
	w.req.get(path, &u)
	return u.Response.User.Uris.UserAlbums.URI
}

// albums make multiple calls to obtain the full list of user albums. It calls the albums endpoint
// unless the "NextPage" value in the response is empty
func (w *Worker) albums(firstURI string) ([]album, error) {
	uri := firstURI
	var albums []album
	for uri != "" {
		var a albumsResponse
		if err := w.req.get(uri, &a); err != nil {
			return albums, fmt.Errorf("error getting albums from %s. Error: %v", uri, err)
		}
		albums = append(albums, a.Response.Album...)
		uri = a.Response.Pages.NextPage
	}
	return albums, nil
}

// albumImages make multiple calls to obtain all images of an album. It calls the album images
// endpoint unless the "NextPage" value in the response is empty
func (w *Worker) albumImages(firstURI string, albumPath string) ([]albumImage, error) {
	uri := firstURI
	var images []albumImage
	for uri != "" {
		if w.quitting {
			return nil, nil
		}
		var a albumImagesResponse
		if err := w.req.get(uri, &a); err != nil {
			return images, fmt.Errorf("error getting album images from %s. Error: %v", uri, err)
		}

		// If the album is empty, a.Response.AlbumImage is missing instead of an empty array (weird...)
		if a.Response.AlbumImage == nil {
			log.Infof("album is empty: %s", albumPath)

			break
		}

		// Loop over response in inject the albumPath and then append to the images
		for _, i := range a.Response.AlbumImage {
			i.AlbumPath = albumPath
			if err := i.buildFilename(w.filenameTmpl); err != nil {
				return nil, fmt.Errorf("cannot build image filename: %v", err)
			}
			images = append(images, i)
		}
		uri = a.Response.Pages.NextPage
	}

	return images, nil
}

func (w *Worker) imageTimestamp(img albumImage) time.Time {
	var i imageMetadataResponse
	if err := w.req.get(img.Uris.ImageMetadata.Uri, &i); err != nil {
		return time.Time{}
	}
	return i.Response.DateTimeCreated
}

// saveImages calls saveImage or saveVideo to save a list of album images to the given folder
func (w *Worker) saveImages(images []albumImage, folder string) {
	for _, image := range images {
		if w.quitting {
			return
		}
		w.downloadsCh <- &downloadInfo{
			image:  image,
			folder: folder,
		}
	}
}

// saveImage saves an image to the given folder unless its name is empty
func (w *Worker) saveImage(image albumImage, folder string) error {
	if image.Name() == "" {
		return errors.New("unable to find valid image filename, skipping")
	}
	dest := fmt.Sprintf("%s/%s", folder, image.Name())
	log.Debug(image.ArchivedUri)

	ok, err := w.downloadFn(dest, image.ArchivedUri, image.ArchivedSize)
	if err != nil {
		return err
	}

	if w.cfg.UseMetadataTimes && (ok || w.cfg.ForceMetadataTimes) {
		return w.setChTime(image, dest)
	}

	return nil
}

// saveVideo saves a video to the given folder unless its name is empty or is still under processing
func (w *Worker) saveVideo(image albumImage, folder string) error {
	if image.Name() == "" {
		return errors.New("unable to find valid video filename, skipping")
	}
	dest := fmt.Sprintf("%s/%s", folder, image.Name())

	if image.Processing && !w.cfg.ForceVideoDownload { // Skip videos if under processing
		return fmt.Errorf("skipping video %s because under processing, %#v", image.Name(), image)
	}

	var v albumVideo
	log.Debug("(saveVideo) getting ", image.Uris.LargestVideo.Uri)
	if err := w.req.get(image.Uris.LargestVideo.Uri, &v); err != nil {
		return fmt.Errorf("cannot get URI for video %+v. Error: %v", image, err)
	}

	ok, err := w.downloadFn(dest, v.Response.LargestVideo.Url, v.Response.LargestVideo.Size)
	if err != nil {
		return err
	}

	if w.cfg.UseMetadataTimes && (ok || w.cfg.ForceMetadataTimes) {
		return w.setChTime(image, dest)
	}

	return nil
}

func (w *Worker) setChTime(image albumImage, dest string) error {
	// Try first with the date in the image, to avoid making an additional call
	created, err := time.Parse(time.RFC3339, image.DateTimeOriginal)
	if err != nil || created.IsZero() {
		created = w.imageTimestamp(image)
	}
	if !created.IsZero() {
		log.Debugf("setting chtime %v for %s", created, dest)
		return os.Chtimes(dest, time.Now(), created)
	}

	return nil
}
