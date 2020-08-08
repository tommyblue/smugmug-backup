package smugmug_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/tommyblue/smugmug-backup"
)

func TestNew(t *testing.T) {
	// Create a real folder
	ok_dir, err := ioutil.TempDir("/tmp", "smugmug-backup")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(ok_dir)

	// Create and immediately delete a folder, so we are sure it doesn't exist
	unexisting_dir, err := ioutil.TempDir("/tmp", "smugmug-backup")
	if err != nil {
		t.Fatal(err)
	}
	os.RemoveAll(unexisting_dir)

	tests := []struct {
		name        string
		username    string
		destination string
		wantErr     bool
	}{
		{
			name:        "empty name",
			username:    "",
			destination: "/path/to/dest/",
			wantErr:     true,
		},
		{
			name:        "empty destination",
			username:    "user_name",
			destination: "",
			wantErr:     true,
		},
		{
			name:        "correct",
			username:    "user_name",
			destination: unexisting_dir,
			wantErr:     true,
		},
		{
			name:        "correct",
			username:    "user_name",
			destination: ok_dir,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := smugmug.New(&smugmug.Conf{
				Username:    tt.username,
				Destination: tt.destination,
			})
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantErr && err == nil {
				t.Fatalf("want error, got nil")
			}
		})
	}

	// if err := s.Run(); err != nil {
	// 	log.Fatal(err)
	// }
}
