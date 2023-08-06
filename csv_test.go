package smugmug

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func Test_createMetadataCSV(t *testing.T) {
	fpath := filepath.Join(t.TempDir(), "file.csv")
	if err := createMetadataCSV(fpath); err != nil {
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

func Test_writeToCSV(t *testing.T) {
	fpath := filepath.Join(t.TempDir(), "file.csv")
	if err := createMetadataCSV(fpath); err != nil {
		t.Fatalf("cannot create csv file: %v", err)
	}

	w := &Worker{
		cfg: &Conf{
			metadataFile: fpath,
		},
	}

	images := []albumImage{
		{
			builtFilename: "fname1",
			ArchivedUri:   "url",
			Caption:       "asdsad",
			Keywords:      "a,b,c",
			Latitude:      "40.123",
			Longitude:     "11.11",
		},
		{
			builtFilename: "fname2",
			ArchivedUri:   "url",
			Caption:       "asdsad",
			Keywords:      "a,b,c",
			Latitude:      "40.123",
			Longitude:     "11.11",
		},
		{
			builtFilename: "fname3",
			ArchivedUri:   "url",
			Caption:       "asdsad",
			Keywords:      "a,b,c",
			Latitude:      "40.123",
			Longitude:     "11.11",
		},
	}

	w.writeToCSV(images)

	f, err := os.Open(fpath)
	if err != nil {
		t.Fatalf("cannot open csv file: %v", err)
	}

	n, err := lineCounter(t, f)
	if err != nil {
		t.Fatalf("cannot count lines in csv file: %v", err)
	}

	if n != 4 {
		t.Fatalf("want 4 lines, got %d", n)
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
