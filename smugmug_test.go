package main

import (
	"testing"
)

type mockHandler struct {
	called int
}

func (c *mockHandler) get(url string, obj interface{}) error {
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

func Test_smugMugConf_getAlbums(t *testing.T) {
	c := &smugMugConf{
		r: &mockHandler{},
	}
	u := &user{}
	u.Response.User.Uris.UserAlbums.URI = "someurl"
	albums, err := c.getAlbums(u)
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	if c.r.(*mockHandler).called != 2 {
		t.Errorf("Called, want 2, got %d", c.r.(*mockHandler).called)
	}
	if len(*albums) != 4 {
		t.Errorf("Want 4, got %d", len(*albums))
	}
}
