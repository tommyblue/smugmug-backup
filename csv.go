package smugmug

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// ALBUMS_METADATA_FNAME is the name of the CSV file used to store albums metadata
const ALBUMS_METADATA_FNAME = "albums_metadata.csv"

// IMAGES_METADATA_FNAME is the name of the CSV file used to store images/videos metadata
const IMAGES_METADATA_FNAME = "images_metadata.csv"

var albumsCsvHeader = []string{
	"AlbumKey",
	"URLPath",
	"Title",
	"Name",
	"NiceName",
	"Description",
	"Keywords",
	"Date",
	"LastUpdated",
	"ImagesLastUpdated",
	"ImageCount",
	"Privacy",
	"SecurityType",
	"SortMethod",
	"SortDirection",
	"WebUri",
	"AllowDownloads",
	"PasswordHint",
	"Protected",
	"HighlightImageUri",
	"HighlightImageKey",
}

var imagesCsvHeader = []string{
	"Filename",
	"Type",
	"AlbumKey",
	"AlbumPath",
	"IsHighlight",
	"ImageKey",
	"Title",
	"Caption",
	"Keywords",
	"Format",
	"Width",
	"Height",
	"OriginalWidth",
	"OriginalHeight",
	"Size",
	"ArchivedSize",
	"ArchivedUri",
	"ArchivedMD5",
	"DateTimeOriginal",
	"DateTimeUploaded",
	"Latitude",
	"Longitude",
	"Altitude",
	"Hidden",
	"Watermarked",
	"Collectable",
	"IsArchive",
	"Status",
	"SubStatus",
	"WebUri",
	"ThumbnailUrl",
	"UploadKey",
}

// createAlbumsMetadataCSV creates the albums metadata CSV file and writes the header line
func createAlbumsMetadataCSV(fpath string) error {
	file, err := os.Create(fpath)
	if err != nil {
		return fmt.Errorf("failed creating albums metadata CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	return writer.Write(albumsCsvHeader)
}

// createImagesMetadataCSV creates the images metadata CSV file and writes the header line
func createImagesMetadataCSV(fpath string) error {
	file, err := os.Create(fpath)
	if err != nil {
		return fmt.Errorf("failed creating images metadata CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	return writer.Write(imagesCsvHeader)
}

// buildAlbumMetadata returns the data to be added to the albums metadata CSV file
func (w *Worker) buildAlbumMetadata(a album) []string {
	return []string{
		a.AlbumKey,
		a.URLPath,
		a.Title,
		a.Name,
		a.NiceName,
		a.Description,
		a.Keywords,
		a.Date,
		a.LastUpdated,
		a.ImagesLastUpdated,
		strconv.Itoa(a.ImageCount),
		a.Privacy,
		a.SecurityType,
		a.SortMethod,
		a.SortDirection,
		a.WebUri,
		strconv.FormatBool(a.AllowDownloads),
		a.PasswordHint,
		strconv.FormatBool(a.Protected),
		a.HighlightImageUri(),
		a.HighlightImageKey(),
	}
}

// buildImageMetadata returns the data to be added to the images metadata CSV file
func (w *Worker) buildImageMetadata(a albumImage, alb album, folder string) []string {
	ftype := "image"
	if a.IsVideo {
		ftype = "video"
	}
	
	// Check if this image is the album highlight
	isHighlight := a.ImageKey == alb.HighlightImageKey()

	return []string{
		fmt.Sprintf("%s/%s", folder, a.Name()),
		ftype,
		a.AlbumKey,
		a.AlbumPath,
		strconv.FormatBool(isHighlight),
		a.ImageKey,
		a.Title,
		a.Caption,
		a.Keywords,
		a.Format,
		strconv.Itoa(a.Width),
		strconv.Itoa(a.Height),
		strconv.Itoa(a.OriginalWidth),
		strconv.Itoa(a.OriginalHeight),
		strconv.FormatInt(a.Size, 10),
		strconv.FormatInt(a.ArchivedSize, 10),
		a.ArchivedUri,
		a.ArchivedMD5,
		a.DateTimeOriginal,
		a.DateTimeUploaded,
		a.Latitude,
		a.Longitude,
		strconv.Itoa(a.Altitude),
		strconv.FormatBool(a.Hidden),
		strconv.FormatBool(a.Watermarked),
		strconv.FormatBool(a.Collectable),
		strconv.FormatBool(a.IsArchive),
		a.Status,
		a.SubStatus,
		a.WebUri,
		a.ThumbnailUrl,
		a.UploadKey,
	}
}

// writeAlbumToCSV writes album metadata to CSV file
func (w *Worker) writeAlbumToCSV(alb album) {
	w.csvLock.Lock()
	defer w.csvLock.Unlock()

	file, err := os.OpenFile(w.cfg.albumsMetadataFile, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Errorf("cannot open albums metadata CSV file: %v", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer func() {
		writer.Flush()
		if err := writer.Error(); err != nil {
			log.Errorf("error flushing albums metadata file: %v", err)
		}
	}()

	if err := writer.Write(w.buildAlbumMetadata(alb)); err != nil {
		log.Errorf("cannot write to albums metadata CSV file: %v", err)
	}
}

// writeImagesToCSV writes images metadata to CSV file
func (w *Worker) writeImagesToCSV(images []albumImage, alb album, folder string) {
	w.csvLock.Lock()
	defer w.csvLock.Unlock()

	file, err := os.OpenFile(w.cfg.imagesMetadataFile, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Errorf("cannot open images metadata CSV file: %v", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer func() {
		writer.Flush()
		if err := writer.Error(); err != nil {
			log.Errorf("error flushing images metadata file: %v", err)
		}
	}()

	for _, img := range images {
		if err := writer.Write(w.buildImageMetadata(img, alb, folder)); err != nil {
			log.Errorf("cannot write to images metadata CSV file: %v", err)
		}
	}
}
