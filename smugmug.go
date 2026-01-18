package smugmug

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Conf is the configuration of the smugmug worker
type Conf struct {
	ApiKey              string // API key
	ApiSecret           string // API secret
	UserToken           string // User token
	UserSecret          string // User secret
	Destination         string // Backup destination folder
	Filenames           string // Template for files naming
	UseMetadataTimes    bool   // When true, the last update timestamp will be retrieved from metadata
	ForceMetadataTimes  bool   // When true, then the last update timestamp is always retrieved and overwritten, also for existing files
	WriteCSV            bool   // When true, a CSV file including downloaded files metadata is written
	ForceVideoDownload  bool   // When true, download videos also if marked as under processing
	ConcurrentDownloads int    // number of concurrent downloads of images and videos, default is 1
	ConcurrentAlbums    int    // number of concurrent albums analyzed via API calls
	HTTPBaseUrl         string // Smugmug API URL, defaults to https://api.smugmug.com
	HTTPMaxRetries      int    // Max number of retries for HTTP calls, defaults to 3

	username            string
	albumsMetadataFile  string
	imagesMetadataFile  string
}

// overrideEnvConf overrides any configuration value if the
// corresponding environment variables is set
func (cfg *Conf) overrideEnvConf() {
	if os.Getenv("SMGMG_BK_API_KEY") != "" {
		cfg.ApiKey = os.Getenv("SMGMG_BK_API_KEY")
	}

	if os.Getenv("SMGMG_BK_API_SECRET") != "" {
		cfg.ApiSecret = os.Getenv("SMGMG_BK_API_SECRET")
	}

	if os.Getenv("SMGMG_BK_USER_TOKEN") != "" {
		cfg.UserToken = os.Getenv("SMGMG_BK_USER_TOKEN")
	}

	if os.Getenv("SMGMG_BK_USER_SECRET") != "" {
		cfg.UserSecret = os.Getenv("SMGMG_BK_USER_SECRET")
	}

	if os.Getenv("SMGMG_BK_DESTINATION") != "" {
		cfg.Destination = os.Getenv("SMGMG_BK_DESTINATION")
	}

	if os.Getenv("SMGMG_BK_FILE_NAMES") != "" {
		cfg.Filenames = os.Getenv("SMGMG_BK_FILE_NAMES")
	}
}

func (cfg *Conf) validate() error {
	if cfg.ApiKey == "" {
		return errors.New("ApiKey can't be empty")
	}

	if cfg.ApiSecret == "" {
		return errors.New("ApiSecret can't be empty")
	}

	if cfg.UserToken == "" {
		return errors.New("UserToken can't be empty")
	}

	if cfg.UserSecret == "" {
		return errors.New("UserSecret can't be empty")
	}

	if cfg.Destination == "" {
		return errors.New("destination can't be empty")
	}

	// Check exising and writeability of destination folder
	if err := checkDestFolder(cfg.Destination); err != nil {
		return fmt.Errorf("can't find in the destination folder %s: %v", cfg.Destination, err)
	}

	return nil
}

// ReadConf produces a configuration object for the Smugmug worker.
//
// It reads the configuration from ./config.toml or "$HOME/.smgmg/config.toml"
func ReadConf(cfgPath string) (*Conf, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	if cfgPath != "" {
		viper.AddConfigPath(cfgPath)
	}
	viper.AddConfigPath("$HOME/.smgmg")
	viper.AddConfigPath(".")

	// defaults
	viper.SetDefault("http.base_url", "https://api.smugmug.com")
	viper.SetDefault("http.max_retries", 3)
	viper.SetDefault("store.file_names", "{{.FileName}}")
	viper.SetDefault("store.concurrent_downloads", 1)
	viper.SetDefault("store.concurrent_albums", 1)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("configuration file not found in ./config.toml or $HOME/.smgmg/config.toml")
		} else {
			return nil, err
		}
	}

	if viper.GetString("authentication.username") != "" {
		log.Warnf("[DEPRECATION] Username configuration value is ignored. It is now retrieved automatically from SmugMug based on the authentication credentials.")
	}

	cfg := &Conf{
		ApiKey:              viper.GetString("authentication.api_key"),
		ApiSecret:           viper.GetString("authentication.api_secret"),
		UserToken:           viper.GetString("authentication.user_token"),
		UserSecret:          viper.GetString("authentication.user_secret"),
		Destination:         viper.GetString("store.destination"),
		Filenames:           viper.GetString("store.file_names"),
		UseMetadataTimes:    viper.GetBool("store.use_metadata_times"),
		ForceMetadataTimes:  viper.GetBool("store.force_metadata_times"),
		WriteCSV:            viper.GetBool("store.write_csv"),
		ForceVideoDownload:  viper.GetBool("store.force_video_download"),
		ConcurrentDownloads: viper.GetInt("store.concurrent_downloads"),
		ConcurrentAlbums:    viper.GetInt("store.concurrent_albums"),
		HTTPBaseUrl:         viper.GetString("http.base_url"),
		HTTPMaxRetries:      viper.GetInt("http.max_retries"),
	}

	cfg.overrideEnvConf()

	if !cfg.UseMetadataTimes && cfg.ForceMetadataTimes {
		return nil, errors.New("cannot use store.force_metadata_times without store.use_metadata_times")
	}

	return cfg, nil
}

type FileMetadata struct {
	FileName    string
	ArchivedUri string
	Caption     string
	Keywords    string
}

type downloadInfo struct {
	image  albumImage
	folder string
}

// Worker actually implements the backup logic
type Worker struct {
	req              requestsHandler
	cfg              *Conf
	errors           int
	downloadFn       func(string, string, int64) (bool, error) // defined in struct for better testing
	filenameTmpl     *template.Template
	downloadsCh      chan *downloadInfo
	downloadsWorkers int
	downloadWg       sync.WaitGroup
	stopCh           chan struct{}
	quitting         bool
	albumCh          chan album
	albumsWorkers    int
	albumWg          sync.WaitGroup
	csvLock          sync.Mutex
}

// New return a SmugMug backup configuration. It returns an error if it fails parsing
// the command line arguments
func New(cfg *Conf) (*Worker, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	handler := newHTTPHandler(cfg.HTTPBaseUrl, cfg.HTTPMaxRetries, cfg.ApiKey, cfg.ApiSecret, cfg.UserToken, cfg.UserSecret)

	tmpl, err := buildFilenameTemplate(cfg.Filenames)
	if err != nil {
		return nil, err
	}

	if cfg.WriteCSV {
		cfg.albumsMetadataFile = filepath.Join(cfg.Destination, ALBUMS_METADATA_FNAME)
		cfg.imagesMetadataFile = filepath.Join(cfg.Destination, IMAGES_METADATA_FNAME)
		createAlbumsMetadataCSV(cfg.albumsMetadataFile)
		createImagesMetadataCSV(cfg.imagesMetadataFile)
	}

	return &Worker{
		cfg:              cfg,
		req:              handler,
		downloadFn:       handler.download,
		filenameTmpl:     tmpl,
		downloadsCh:      make(chan *downloadInfo),
		downloadsWorkers: cfg.ConcurrentDownloads,
		downloadWg:       sync.WaitGroup{},
		stopCh:           make(chan struct{}),
		albumCh:          make(chan album),
		albumsWorkers:    cfg.ConcurrentAlbums,
		albumWg:          sync.WaitGroup{},
	}, nil
}

func (w *Worker) albumWorker(id int) {
	log.Debugf("Running albumWorker %d", id)
	for {
		select {
		case <-w.stopCh:
			log.Debugf("Stopping albumWorker %d", id)
			return
		case album, ok := <-w.albumCh:
			if !ok {
				// Channel is closed
				log.Debugf("Quitting albumWorker %d", id)
				return
			}
			folder := filepath.Join(w.cfg.Destination, album.URLPath)

			if err := createFolder(folder); err != nil {
				log.WithError(err).Errorf("cannot create the destination folder %s", folder)
				w.errors++
				continue
			}

			log.Debugf("[ALBUM IMAGES] %s", album.Uris.AlbumImages.URI)
			images, err := w.albumImages(album.Uris.AlbumImages.URI, album)
			if err != nil {
				log.WithError(err).Errorf("cannot get album images for %s", album.Uris.AlbumImages.URI)
				w.errors++
				continue
			}

			log.Debugf("Got album images for %s", album.Uris.AlbumImages.URI)
			// log.Debugf("%+v", images)
			w.saveImages(images, folder)
			if w.cfg.WriteCSV {
				w.writeAlbumToCSV(album)
				w.writeImagesToCSV(images, folder)
			}
		}
	}
}

func (w *Worker) downloader(id int) {
	log.Debugf("Running downloader %d", id)
	for {
		select {
		case <-w.stopCh:
			log.Debugf("Stopping downloader %d", id)
			return
		case info, ok := <-w.downloadsCh:
			if !ok {
				log.Debugf("Quitting downloader %d", id)
				return
			}

			if info.image.IsVideo {
				if err := w.saveVideo(info.image, info.folder); err != nil {
					log.Warnf("Error: %v", err)
				}
				continue
			}

			if err := w.saveImage(info.image, info.folder); err != nil {
				log.Warnf("Error: %v", err)
			}
		}
	}
}

func buildFilenameTemplate(filenameTemplate string) (*template.Template, error) {
	// Use FileName as default
	if filenameTemplate == "" {
		filenameTemplate = "{{.FileName}}"
	}
	tmpl, err := template.New("image_filename").Option("missingkey=error").Parse(filenameTemplate)
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

// Run performs the backup of the provided SmugMug account.
//
// The workflow is the following:
//
//   - Get user albums
//   - Iterate over all albums and:
//   - create folder
//   - iterate over all images and videos
//   - if existing and with the same size, then skip
//   - if not, download
func (w *Worker) Run() error {
	var err error
	w.cfg.username, err = w.currentUser()
	if err != nil {
		return fmt.Errorf("error checking credentials: %v", err)
	}

	w.albumWg.Add(w.albumsWorkers)
	for i := 0; i < w.albumsWorkers; i++ {
		go func(i int) {
			defer w.albumWg.Done()
			w.albumWorker(i)
		}(i)
	}

	w.downloadWg.Add(w.downloadsWorkers)
	for i := 0; i < w.downloadsWorkers; i++ {
		go func(i int) {
			defer w.downloadWg.Done()
			w.downloader(i)
		}(i)
	}

	// Get user albums
	log.Infof("Getting albums for user %s...\n", w.cfg.username)
	albums, err := w.userAlbums()
	if err != nil {
		return fmt.Errorf("error getting user albums: %v", err)
	}

	log.Infof("Found %d albums\n", len(albums))

	for _, album := range albums {
		if w.quitting {
			break
		}
		w.albumCh <- album
	}

	w.Wait()

	if w.errors > 0 {
		return fmt.Errorf("completed with %d errors, please check logs", w.errors)
	}

	if w.quitting {
		log.Info("Quit worker!")
		return nil
	}

	log.Info("Backup completed.")
	return nil
}

func (w *Worker) Stop() {
	log.Info("Quitting worker...")
	close(w.stopCh)
	w.quitting = true
}

func (w *Worker) Wait() {
	close(w.albumCh)
	log.Debug("waiting albumWg...")
	w.albumWg.Wait()
	log.Debug("albumWg done.")

	close(w.downloadsCh)
	log.Debug("waiting downloadWg...")
	w.downloadWg.Wait()
	log.Debug("downloadWg done.")
}
