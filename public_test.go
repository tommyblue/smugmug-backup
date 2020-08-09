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
		apiKey      string
		apiSecret   string
		userToken   string
		userSecret  string
		wantErr     bool
	}{
		{
			name:        "empty name",
			username:    "",
			destination: "/path/to/dest/",
			apiKey:      "value",
			apiSecret:   "value",
			userToken:   "value",
			userSecret:  "value",
			wantErr:     true,
		},
		{
			name:        "empty destination",
			username:    "user_name",
			destination: "",
			apiKey:      "value",
			apiSecret:   "value",
			userToken:   "value",
			userSecret:  "value",
			wantErr:     true,
		},
		{
			name:        "unexisting dest dir",
			username:    "user_name",
			destination: unexisting_dir,
			apiKey:      "value",
			apiSecret:   "value",
			userToken:   "value",
			userSecret:  "value",
			wantErr:     true,
		},
		{
			name:        "missing apiKey",
			username:    "user_name",
			destination: ok_dir,
			apiKey:      "",
			apiSecret:   "value",
			userToken:   "value",
			userSecret:  "value",
			wantErr:     true,
		},
		{
			name:        "missing apiSecret",
			username:    "user_name",
			destination: ok_dir,
			apiKey:      "value",
			apiSecret:   "",
			userToken:   "value",
			userSecret:  "value",
			wantErr:     true,
		},
		{
			name:        "missing userToken",
			username:    "user_name",
			destination: ok_dir,
			apiKey:      "value",
			apiSecret:   "value",
			userToken:   "",
			userSecret:  "value",
			wantErr:     true,
		},
		{
			name:        "missing userSecret",
			username:    "user_name",
			destination: ok_dir,
			apiKey:      "value",
			apiSecret:   "value",
			userToken:   "value",
			userSecret:  "",
			wantErr:     true,
		},
		{
			name:        "correct",
			username:    "user_name",
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
				Username:    tt.username,
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
