package smugmug

import (
	"bytes"
	"encoding/csv"
	"io"
	"os"
	"path/filepath"
	"testing"
)

// mockHighlightHandler mocks the request handler for highlight node API calls
type mockHighlightHandler struct {
	highlightImageKey string
}

func (m *mockHighlightHandler) get(url string, obj interface{}) error {
	// Return mock highlight response
	if resp, ok := obj.(*highlightNodeResponse); ok {
		resp.Response.Image.ImageKey = m.highlightImageKey
		return nil
	}
	return nil
}

func (m *mockHighlightHandler) download(dest, downloadURL string, fileSize int64) (bool, error) {
	return false, nil
}

func Test_createAlbumsMetadataCSV(t *testing.T) {
	fpath := filepath.Join(t.TempDir(), "albums.csv")
	if err := createAlbumsMetadataCSV(fpath); err != nil {
		t.Fatalf("cannot create csv file: %v", err)
	}

	f, err := os.Open(fpath)
	if err != nil {
		t.Fatalf("cannot open csv file: %v", err)
	}

	n, err := lineCounter(t, f)
	if err != nil {
		t.Fatalf("cannot count lines in csv file: %v", err)
	}

	if n != 1 {
		t.Fatalf("want 1 line, got %d", n)
	}
}

func Test_createImagesMetadataCSV(t *testing.T) {
	fpath := filepath.Join(t.TempDir(), "images.csv")
	if err := createImagesMetadataCSV(fpath); err != nil {
		t.Fatalf("cannot create csv file: %v", err)
	}

	f, err := os.Open(fpath)
	if err != nil {
		t.Fatalf("cannot open csv file: %v", err)
	}

	n, err := lineCounter(t, f)
	if err != nil {
		t.Fatalf("cannot count lines in csv file: %v", err)
	}

	if n != 1 {
		t.Fatalf("want 1 line, got %d", n)
	}
}

func Test_writeImagesToCSV(t *testing.T) {
	fpath := filepath.Join(t.TempDir(), "images.csv")
	if err := createImagesMetadataCSV(fpath); err != nil {
		t.Fatalf("cannot create csv file: %v", err)
	}

	// Mock handler that returns image1 as the highlight ImageKey
	mockHandler := &mockHighlightHandler{highlightImageKey: "image1"}

	w := &Worker{
		cfg: &Conf{
			imagesMetadataFile: fpath,
		},
		req: mockHandler,
	}

	alb := album{
		AlbumKey: "testAlbum",
	}
	// Set the highlight image URI to a node URI
	alb.Uris.HighlightImage.URI = "/api/v2/highlight/node/testNodeId"

	images := []albumImage{
		{
			builtFilename: "fname1",
			ArchivedUri:   "url",
			Caption:       "first image",
			Keywords:      "a,b,c",
			Latitude:      "40.123",
			Longitude:     "11.11",
			AlbumKey:      "album1",
			ImageKey:      "image1", // This is the highlight image
		},
		{
			builtFilename: "fname2",
			ArchivedUri:   "url",
			Caption:       "second image",
			Keywords:      "a,b,c",
			Latitude:      "40.123",
			Longitude:     "11.11",
			AlbumKey:      "album1",
			ImageKey:      "image2",
		},
		{
			builtFilename: "fname3",
			ArchivedUri:   "url",
			Caption:       "third image",
			Keywords:      "a,b,c",
			Latitude:      "40.123",
			Longitude:     "11.11",
			AlbumKey:      "album1",
			ImageKey:      "image3",
		},
	}

	w.writeImagesToCSV(images, alb, "test")

	f, err := os.Open(fpath)
	if err != nil {
		t.Fatalf("cannot open csv file: %v", err)
	}
	defer f.Close()

	// Verify structure
	reader := csv.NewReader(f)
	
	// Read header
	header, err := reader.Read()
	if err != nil {
		t.Fatalf("cannot read header: %v", err)
	}
	
	// Find IsHighlight column
	isHighlightIdx := -1
	imageKeyIdx := -1
	for i, col := range header {
		if col == "IsHighlight" {
			isHighlightIdx = i
		}
		if col == "ImageKey" {
			imageKeyIdx = i
		}
	}
	
	if isHighlightIdx == -1 {
		t.Fatalf("IsHighlight column not found in header")
	}
	if imageKeyIdx == -1 {
		t.Fatalf("ImageKey column not found in header")
	}

	// Read and verify rows
	rowCount := 0
	highlightImageKey := ""
	
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("error reading row: %v", err)
		}
		
		rowCount++
		// Record the image key that has IsHighlight=true
		if record[isHighlightIdx] == "true" {
			highlightImageKey = record[imageKeyIdx]
		}
	}

	if rowCount != 3 {
		t.Fatalf("want 3 data rows, got %d", rowCount)
	}
	
	// Verify that image1 is marked as highlight
	if highlightImageKey != "image1" {
		t.Fatalf("expected image1 to be highlight, got %s", highlightImageKey)
	}
}

func lineCounter(t *testing.T, r io.Reader) (int, error) {
	t.Helper()
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
