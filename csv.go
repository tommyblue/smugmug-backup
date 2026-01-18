package smugmug

import (
	"encoding/csv"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

// METADATA_FNAME is the name of the CSV file used to store files metadata
const METADATA_FNAME = "metadata.csv"

var csvHeader = []string{
	"Filename",
	"Type",
	"ArchivedUri",
	"Caption",
	"Keywords",
	"Latitude",
	"Longitude",
	"AlbumTitle",
	"AlbumDescription",
	"AlbumKeywords",
	"AlbumCreated",
	"AlbumLastUpdated",
	"DateTimeOriginal",
	"DateTimeUploaded",
}

// createMetadataCSV creates the metadata CSV file and writes the header line
func createMetadataCSV(fpath string) error {
	file, err := os.Create(fpath)
	if err != nil {
		return fmt.Errorf("failed creating metadata CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	return writer.Write(csvHeader)
}

// buildMetadata returns the data to be added to the metadata CSV file
func (w *Worker) buildMetadata(a albumImage, folder string) []string {
	ftype := "image"
	if a.IsVideo {
		ftype = "video"
	}

	return []string{
		fmt.Sprintf("%s/%s", folder, a.Name()),
		ftype,
		a.ArchivedUri,
		a.Caption,
		a.Keywords,
		a.Latitude,
		a.Longitude,
		a.AlbumTitle,
		a.AlbumDescription,
		a.AlbumKeywords,
		a.AlbumCreated,
		a.AlbumLastUpdated,
		a.DateTimeOriginal,
		a.DateTimeUploaded,
	}
}

// writeToCSV writes images metadata to CSV file
func (w *Worker) writeToCSV(images []albumImage, folder string) {
	w.csvLock.Lock()
	defer w.csvLock.Unlock()

	file, err := os.OpenFile(w.cfg.metadataFile, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Errorf("cannot open metadata CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer func() {
		writer.Flush()
		if err := writer.Error(); err != nil {
			log.Errorf("error flushing metadata file: %v", err)
		}
	}()

	for _, img := range images {
		if err := writer.Write(w.buildMetadata(img, folder)); err != nil {
			log.Errorf("cannot write to metadata CSV file: %v", err)
		}
	}
}
