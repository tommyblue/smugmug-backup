package main

import (
	"fmt"
)

type node struct{}

type requestsHandler interface {
	get(string, interface{}) error
}

type smugmugHandler struct{}

func (c *smugMugConf) getUser() *user {
	var u user
	path := fmt.Sprintf("/api/v2/user/%s", c.username)
	c.r.get(path, &u)
	return &u
}

func (c *smugMugConf) getAlbums(u *user) (*[]album, error) {
	var albums []album
	url := u.Response.User.Uris.UserAlbums.URI
	for url != "" {
		var a albumsResponse
		c.r.get(url, &a)
		albums = append(albums, a.Response.Album...)
		url = a.Response.Pages.NextPage
	}
	return &albums, nil
}

func (c *smugMugConf) getAlbumImages(ImagesURL string) (*[]albumImage, error) {
	var albumImages albumImagesResponse
	var images []albumImage
	c.r.get(ImagesURL, &albumImages)
	images = append(images, albumImages.Response.AlbumImage...)
	nextPage := albumImages.Response.Pages.NextPage
	for nextPage != "" {
		var albumImages albumImagesResponse
		c.r.get(nextPage, &albumImages)
		nextPage = albumImages.Response.Pages.NextPage
		images = append(images, albumImages.Response.AlbumImage...)
	}
	return &images, nil
}
