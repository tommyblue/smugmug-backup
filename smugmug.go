package main

import (
	"fmt"
)

type node struct{}

func (c *smugMugConf) getAlbums() (*[]album, error) {
	var u user
	path := fmt.Sprintf("/api/v2/user/%s", c.username)
	c.get(path, &u)

	var albums []album
	url := u.Response.User.Uris.UserAlbums.URI
	for url != "" {
		var a albumsResponse
		c.get(url, &a)
		albums = append(albums, a.Response.Album...)
		url = a.Response.Pages.NextPage
	}
	return &albums, nil
}

func (c *smugMugConf) getAlbumImages(ImagesURL string) (*[]albumImage, error) {
	var albumImages albumImagesResponse
	var images []albumImage
	c.get(ImagesURL, &albumImages)
	images = append(images, albumImages.Response.AlbumImage...)
	nextPage := albumImages.Response.Pages.NextPage
	for nextPage != "" {
		var albumImages albumImagesResponse
		c.get(nextPage, &albumImages)
		nextPage = albumImages.Response.Pages.NextPage
		images = append(images, albumImages.Response.AlbumImage...)
	}
	return &images, nil
}
