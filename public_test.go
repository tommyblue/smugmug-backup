package smugmug_test

import (
	"os"
	"testing"

	"github.com/tommyblue/smugmug-backup"
)

func TestNew(t *testing.T) {
	// Create a real folder
	ok_dir := t.TempDir()

	// Create and immediately delete a folder, so we are sure it doesn't exist
	unexisting_dir := t.TempDir()
	os.RemoveAll(unexisting_dir)

	tests := []struct {
		name        string
		destination string
		apiKey      string
		apiSecret   string
		userToken   string
		userSecret  string
		wantErr     bool
	}{
		{
			name:        "empty name",
			destination: "/path/to/dest/",
			apiKey:      "value",
			apiSecret:   "value",
			userToken:   "value",
			userSecret:  "value",
			wantErr:     true,
		},
		{
			name:        "empty destination",
			destination: "",
			apiKey:      "value",
			apiSecret:   "value",
			userToken:   "value",
			userSecret:  "value",
			wantErr:     true,
		},
		{
			name:        "unexisting dest dir",
			destination: unexisting_dir,
			apiKey:      "value",
			apiSecret:   "value",
			userToken:   "value",
			userSecret:  "value",
			wantErr:     true,
		},
		{
			name:        "missing apiKey",
			destination: ok_dir,
			apiKey:      "",
			apiSecret:   "value",
			userToken:   "value",
			userSecret:  "value",
			wantErr:     true,
		},
		{
			name:        "missing apiSecret",
			destination: ok_dir,
			apiKey:      "value",
			apiSecret:   "",
			userToken:   "value",
			userSecret:  "value",
			wantErr:     true,
		},
		{
			name:        "missing userToken",
			destination: ok_dir,
			apiKey:      "value",
			apiSecret:   "value",
			userToken:   "",
			userSecret:  "value",
			wantErr:     true,
		},
		{
			name:        "missing userSecret",
			destination: ok_dir,
			apiKey:      "value",
			apiSecret:   "value",
			userToken:   "value",
			userSecret:  "",
			wantErr:     true,
		},
		{
			name:        "correct",
			destination: ok_dir,
			apiKey:      "value",
			apiSecret:   "value",
			userToken:   "value",
			userSecret:  "value",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := smugmug.New(&smugmug.Conf{
				Destination: tt.destination,
				ApiKey:      tt.apiKey,
				ApiSecret:   tt.apiSecret,
				UserToken:   tt.userToken,
				UserSecret:  tt.userSecret,
			})
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantErr && err == nil {
				t.Fatalf("want error, got nil")
			}
		})
	}
}
