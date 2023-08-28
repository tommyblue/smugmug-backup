package smugmug

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
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
	case "/api/v2!authuser":
		var u *currentUser
		u = obj.(*currentUser)
		u.Response.User.NickName = m.username
		return nil
	// from w.userAlbumsURI()
	case fmt.Sprintf("/api/v2/user/%s", testUsername):
		var u *user
		u = obj.(*user)
		u.Response.User.Uris.UserAlbums.URI = m.userAlbumsURI
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

	dest_dir := t.TempDir()

	var downloadCalled int
	tmpl, _ := buildFilenameTemplate("")
	w := &Worker{
		cfg: &Conf{
			Destination: dest_dir,
		},
		req: &mockHandler{
			username:       testUsername,
			userAlbumsURI:  userAlbumsURI,
			albumURLPath:   albumURLPath,
			albumImagesURI: albumImagesURI,
		},
		downloadFn: func(_, _ string, _ int64) (bool, error) {
			downloadCalled++
			return true, nil
		},
		filenameTmpl: tmpl,
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

type testConf struct {
	destination string
	apiKey      string
	apiSecret   string
	userToken   string
	userSecret  string
}

func setupConfFile(t *testing.T, cfg *testConf, missingSecrets bool) func() {
	t.Helper()

	var cfgFile []byte

	if !missingSecrets {
		cfgFile = []byte(fmt.Sprintf(`
[authentication]
api_key = "%s"
api_secret = "%s"
user_token = "%s"
user_secret = "%s"
[store]
destination = "%s"
`, cfg.apiKey, cfg.apiSecret, cfg.userToken, cfg.userSecret, cfg.destination))
	} else {
		cfgFile = []byte(fmt.Sprintf(`
[authentication]
api_key = "%s"
user_token = "%s"
[store]
destination = "%s"
`, cfg.apiKey, cfg.userToken, cfg.destination))
	}

	dir := t.TempDir()

	file := filepath.Join(dir, "config.toml")
	err := os.WriteFile(file, cfgFile, 0644)
	if err != nil {
		t.Fatalf("Can't create config file for testing")
	}

	// push the temp dir to viper to make it find the conf file
	viper.AddConfigPath(dir)

	return func() {
		if err := os.Remove(file); err != nil {
			t.Fatalf("error removing file: %v", err)
		}
		if err := os.Remove(dir); err != nil {
			t.Fatalf("error removing dir: %v", err)
		}
	}
}

func TestReadConf(t *testing.T) {
	viper.Reset()
	cfgObj := &testConf{
		destination: "test_dest",
		apiKey:      "test_apikey",
		apiSecret:   "test_apisecret",
		userToken:   "test_usertoken",
		userSecret:  "test_usersecret",
	}

	defer setupConfFile(t, cfgObj, false)()

	cfg, err := ReadConf("")
	if err != nil {
		t.Fatalf("unexpected err %v", err)
	}

	if cfgObj.destination != cfg.Destination {
		t.Fatalf("destination: want: %s, got: %s", cfgObj.destination, cfg.Destination)
	}

	if cfgObj.apiKey != cfg.ApiKey {
		t.Fatalf("apiKey: want: %s, got: %s", cfgObj.apiKey, cfg.ApiKey)
	}

	if cfgObj.apiSecret != cfg.ApiSecret {
		t.Fatalf("apiSecret: want: %s, got: %s", cfgObj.apiSecret, cfg.ApiSecret)
	}

	if cfgObj.userToken != cfg.UserToken {
		t.Fatalf("userToken: want: %s, got: %s", cfgObj.userToken, cfg.UserToken)
	}

	if cfgObj.userSecret != cfg.UserSecret {
		t.Fatalf("userSecret: want: %s, got: %s", cfgObj.userSecret, cfg.UserSecret)
	}
}

func TestReadConfOverrides(t *testing.T) {
	viper.Reset()
	cfgObj := &testConf{
		destination: "test_dest",
		apiKey:      "test_apikey",
		apiSecret:   "test_apisecret",
		userToken:   "test_usertoken",
		userSecret:  "test_usersecret",
	}

	defer setupConfFile(t, cfgObj, false)()

	os.Setenv("SMGMG_BK_DESTINATION", "overridden_dest")
	os.Setenv("SMGMG_BK_API_KEY", "overridden_apikey")
	os.Setenv("SMGMG_BK_API_SECRET", "overridden_apisecret")
	os.Setenv("SMGMG_BK_USER_TOKEN", "overridden_usertoken")
	os.Setenv("SMGMG_BK_USER_SECRET", "overridden_usersecret")

	defer func() {
		os.Unsetenv("SMGMG_BK_DESTINATION")
		os.Unsetenv("SMGMG_BK_API_KEY")
		os.Unsetenv("SMGMG_BK_API_SECRET")
		os.Unsetenv("SMGMG_BK_USER_TOKEN")
		os.Unsetenv("SMGMG_BK_USER_SECRET")
	}()
	cfg, err := ReadConf("")
	if err != nil {
		t.Fatalf("unexpected err %v", err)
	}

	if cfg.Destination != "overridden_dest" {
		t.Fatalf("destination: want: %s, got: %s", "overridden_dest", cfg.Destination)
	}

	if cfg.ApiKey != "overridden_apikey" {
		t.Fatalf("apiKey: want: %s, got: %s", "overridden_apikey", cfg.ApiKey)
	}

	if cfg.ApiSecret != "overridden_apisecret" {
		t.Fatalf("apiSecret: want: %s, got: %s", "overridden_apisecret", cfg.ApiSecret)
	}

	if cfg.UserToken != "overridden_usertoken" {
		t.Fatalf("userToken: want: %s, got: %s", "overridden_usertoken", cfg.UserToken)
	}

	if cfg.UserSecret != "overridden_usersecret" {
		t.Fatalf("userSecret: want: %s, got: %s", "overridden_usersecret", cfg.UserSecret)
	}
}

func TestReadConfMissingFileValues(t *testing.T) {
	viper.Reset()
	dest_dir := t.TempDir()
	cfgObj := &testConf{
		destination: dest_dir,
		apiKey:      "test_apikey",
		apiSecret:   "test_apisecret",
		userToken:   "test_usertoken",
		userSecret:  "test_usersecret",
	}

	// setup the conf file without the secrets
	defer setupConfFile(t, cfgObj, true)()

	// expected the conf to return an error
	cfg, err := ReadConf("")
	if err != nil {
		t.Fatalf("Unexpected err: %v", err)
	}
	if _, err := New(cfg); err == nil {
		t.Fatalf("expected err")
	}

	// override the secrets
	os.Setenv("SMGMG_BK_API_SECRET", "overridden_apisecret")
	os.Setenv("SMGMG_BK_USER_SECRET", "overridden_usersecret")
	defer func() {
		os.Unsetenv("SMGMG_BK_API_SECRET")
		os.Unsetenv("SMGMG_BK_USER_SECRET")
	}()
	// now it must work
	cfg, err = ReadConf("")
	if err != nil {
		t.Fatalf("Unexpected err: %v", err)
	}
	if _, err := New(cfg); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if cfgObj.destination != cfg.Destination {
		t.Fatalf("destination: want: %s, got: %s", cfgObj.destination, cfg.Destination)
	}

	if cfgObj.apiKey != cfg.ApiKey {
		t.Fatalf("apiKey: want: %s, got: %s", cfgObj.apiKey, cfg.ApiKey)
	}

	if cfgObj.userToken != cfg.UserToken {
		t.Fatalf("userToken: want: %s, got: %s", cfgObj.userToken, cfg.UserToken)
	}

	if cfg.ApiSecret != "overridden_apisecret" {
		t.Fatalf("apiSecret: want: %s, got: %s", "overridden_apisecret", cfg.ApiSecret)
	}

	if cfg.UserSecret != "overridden_usersecret" {
		t.Fatalf("userSecret: want: %s, got: %s", "overridden_usersecret", cfg.UserSecret)
	}
}
