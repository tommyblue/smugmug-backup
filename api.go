package main

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type requestsHandler interface {
	get(string, interface{}) error
}

func (c *smugMugConf) getUserAlbumsURI() string {
	var u user
	path := fmt.Sprintf("/api/v2/user/%s", c.username)
	c.req.get(path, &u)
	return u.Response.User.Uris.UserAlbums.URI
}

func (c *smugMugConf) getAlbums(firstURI string) ([]album, error) {
	uri := firstURI
	var albums []album
	for uri != "" {
		var a albumsResponse
		if err := c.req.get(uri, &a); err != nil {
			return albums, fmt.Errorf("Error getting albums from %s. Error: %v", uri, err)
		}
		albums = append(albums, a.Response.Album...)
		uri = a.Response.Pages.NextPage
	}
	return albums, nil
}

func (c *smugMugConf) getAlbumImages(firstURI string) ([]albumImage, error) {
	uri := firstURI
	var images []albumImage
	for uri != "" {
		var a albumImagesResponse
		if err := c.req.get(uri, &a); err != nil {
			return images, fmt.Errorf("Error getting album images from %s. Error: %v", uri, err)
		}
		images = append(images, a.Response.AlbumImage...)
		uri = a.Response.Pages.NextPage
	}

	return images, nil
}

func (c *smugMugConf) saveImages(images []albumImage, folder string) {
	for _, image := range images {
		if image.IsVideo {
			if err := c.saveVideo(image, folder); err != nil {
				log.Warnf("Error: %v", err)
			}
			continue
		}
		if err := c.saveImage(image, folder); err != nil {
			log.Warnf("Error: %v", err)
		}
	}
}

func (c *smugMugConf) saveImage(image albumImage, folder string) error {
	if image.Name() == "" {
		return errors.New("Unable to find valid image filename, skipping..")
	}
	dest := fmt.Sprintf("%s/%s", folder, image.Name())
	log.Debug(image.ArchivedUri)
	return download(dest, image.ArchivedUri, image.ArchivedSize)
}

func (c *smugMugConf) saveVideo(image albumImage, folder string) error {
	if image.Processing { // Skip videos if under processing
		return fmt.Errorf("Skipping video %s because under processing, %s\n", image.Name(), image.Uris.LargestVideo.Uri)
	}

	if image.Name() == "" {
		return errors.New("Unable to find valid video filename, skipping..")
	}
	dest := fmt.Sprintf("%s/%s", folder, image.Name())

	var v albumVideo
	log.Debug("Getting ", image.Uris.LargestVideo.Uri)
	if err := c.req.get(image.Uris.LargestVideo.Uri, &v); err != nil {
		return fmt.Errorf("Cannot get URI for video %+v. Error: %v", image, err)
	}

	return download(dest, v.Response.LargestVideo.Url, v.Response.LargestVideo.Size)
}
