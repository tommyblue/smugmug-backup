package smugmug

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/tommyblue/smugmug-backup/testutil"
)

const testUsername = "test_username"
const albumURLPath = "dest_path"
const fileName = "filename.jpg"
const userAlbumsURI = "/test_username/albums"
const albumImagesURI = "/album/1/images"

type mockHandler struct {
	username       string
	userAlbumsURI  string
	albumURLPath   string
	albumImagesURI string
}

func (m *mockHandler) get(url string, obj interface{}) error {
	switch url {
	// from w.userAlbumsURI()
	case fmt.Sprintf("/api/v2/user/%s", testUsername):
		var u *user
		u = obj.(*user)
		u.Response.User.Uris.UserAlbums.URI = m.userAlbumsURI
		// log.Infof("DBG: %s", obj.(*user).Response.User.Uris.UserAlbums.URI)
		obj = u
		return nil
	// from w.albums()
	case m.userAlbumsURI:
		albumObj := album{}
		albumObj.URLPath = m.albumURLPath
		albumObj.Uris.AlbumImages.URI = m.albumImagesURI
		albumObjs := []album{albumObj}
		// albumObjs = append(albumObjs, albumObj)
		var a *albumsResponse
		a = obj.(*albumsResponse)
		a.Response.Album = albumObjs
	// from w.albumImages()
	case m.albumImagesURI:
		img1 := albumImage{}
		img1.IsVideo = false
		img1.FileName = fileName
		img1.ImageKey = "abc123"

		img2 := albumImage{}
		img2.IsVideo = false
		img2.FileName = fileName
		img2.ImageKey = "abc124"

		images := []albumImage{img1, img2}

		var a *albumImagesResponse
		a = obj.(*albumImagesResponse)
		a.Response.AlbumImage = images
	}
	return nil
}

func TestRun(t *testing.T) {
	defer testutil.LessLogging()()

	dest_dir, err := ioutil.TempDir("/tmp", "smugmug-backup")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dest_dir)

	var downloadCalled int
	w := &Worker{
		cfg: &Conf{
			Username:    testUsername,
			Destination: dest_dir,
		},
		req: &mockHandler{
			username:       testUsername,
			userAlbumsURI:  userAlbumsURI,
			albumURLPath:   albumURLPath,
			albumImagesURI: albumImagesURI,
		},
		downloadFn: func(_, _ string, _ int64) error {
			downloadCalled++
			return nil
		},
	}
	w.Run()

	dst := filepath.Join(dest_dir, albumURLPath)
	if _, err := os.Stat(dst); err != nil {
		t.Fatalf("Dest folder %s not created", dst)
	}

	if downloadCalled != 2 {
		t.Fatalf("download want %d, got %d", 2, downloadCalled)
	}
}
