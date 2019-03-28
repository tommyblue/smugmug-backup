package main

import (
	"fmt"
)

type node struct{}

func (c *smugMugConf) getAlbums() (*[]album, error) {
	var u user
	path := fmt.Sprintf("/api/v2/user/%s", c.username)
	c.get(path, &u)

	var a albumsResponse
	var albums []album
	c.get(u.Response.User.Uris.UserAlbums.URI, &a)
	albums = append(albums, a.Response.Album...)
	nextPage := a.Response.Pages.NextPage
	for nextPage != "" {
		var a albumsResponse
		c.get(nextPage, &a)
		nextPage = a.Response.Pages.NextPage
		albums = append(albums, a.Response.Album...)
	}
	return &albums, nil
}
