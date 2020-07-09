package main

import (
	"testing"
)

type albumMockHandler struct {
	called int
}

func (c *albumMockHandler) get(url string, obj interface{}) error {
	defer func() { c.called++ }()
	a := obj.(*albumsResponse)
	a.Response.Album = []album{
		{},
		{},
	}
	if c.called == 0 {
		a.Response.Pages.NextPage = "something"
		return nil
	}
	a.Response.Pages.NextPage = ""
	return nil
}

func TestGetAlbums(t *testing.T) {
	c := &smugMugConf{
		req: &albumMockHandler{},
	}
	albums, err := c.getAlbums("someurl")
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	if c.req.(*albumMockHandler).called != 2 {
		t.Errorf("Called, want 2, got %d", c.req.(*albumMockHandler).called)
	}
	if len(albums) != 4 {
		t.Errorf("Want 4, got %d", len(albums))
	}
}

type albumImages struct {
	called int
}

func (c *albumImages) get(url string, obj interface{}) error {
	defer func() { c.called++ }()
	a := obj.(*albumImagesResponse)
	a.Response.AlbumImage = []albumImage{
		{},
		{},
	}
	if c.called == 0 {
		a.Response.Pages.NextPage = "something"
		return nil
	}
	a.Response.Pages.NextPage = ""
	return nil
}

func TestGetAlbumImages(t *testing.T) {
	c := &smugMugConf{
		req: &albumImages{},
	}
	albums, err := c.getAlbumImages("someurl")
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	if c.req.(*albumImages).called != 2 {
		t.Errorf("Called, want 2, got %d", c.req.(*albumImages).called)
	}
	if len(albums) != 4 {
		t.Errorf("Want 4, got %d", len(albums))
	}
}
