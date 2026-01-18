package smugmug

import (
	"bytes"
	"errors"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

type currentUser struct {
	Response struct {
		User struct {
			NickName string `json:"NickName"`
		}
	}
}

type user struct {
	Response struct {
		User struct {
			Uris struct {
				UserAlbums struct {
					URI string `json:"Uri"`
				} `json:"UserAlbums"`
			} `json:"Uris"`
		} `json:"User"`
	} `json:"Response"`
}

type albumsResponse struct {
	Response struct {
		URI   string  `json:"Uri"`
		Album []album `json:"Album"`
		Pages struct {
			NextPage string `json:"NextPage"`
		} `json:"Pages"`
	} `json:"Response"`
}

type album struct {
	URLPath            string `json:"UrlPath"`
	AlbumKey           string `json:"AlbumKey"`
	Title              string `json:"Title"`
	Name               string `json:"Name"`
	NiceName           string `json:"NiceName"`
	Description        string `json:"Description"`
	Keywords           string `json:"Keywords"`
	Date               string `json:"Date"`
	LastUpdated        string `json:"LastUpdated"`
	ImagesLastUpdated  string `json:"ImagesLastUpdated"`
	ImageCount         int    `json:"ImageCount"`
	Privacy            string `json:"Privacy"`
	SecurityType       string `json:"SecurityType"`
	SortMethod         string `json:"SortMethod"`
	SortDirection      string `json:"SortDirection"`
	WebUri             string `json:"WebUri"`
	AllowDownloads     bool   `json:"AllowDownloads"`
	PasswordHint       string `json:"PasswordHint"`
	Protected          bool   `json:"Protected"`
	HighlightImageUri  string `json:"HighlightImageUri"`
	Uris               struct {
		AlbumImages struct {
			URI string `json:"Uri"`
		} `json:"AlbumImages"`
		HighlightImage struct {
			URI string `json:"Uri"`
		} `json:"HighlightImage"`
	} `json:"Uris"`
}

type albumImagesResponse struct {
	Response struct {
		URI        string       `json:"Uri"`
		AlbumImage []albumImage `json:"AlbumImage"`
		Pages      struct {
			NextPage string `json:"NextPage"`
		} `json:"Pages"`
	} `json:"Response"`
}

type imageMetadataResponse struct {
	Response struct {
		DateTimeCreated  time.Time `json:"DateTimeCreated"`
		DateTimeModified time.Time `json:"DateTimeModified"`
	} `json:"Response"`
}

type albumImage struct {
	AlbumPath        string // From album.URLPath
	AlbumKey         string // From album.AlbumKey
	FileName         string `json:"FileName"`
	Title            string `json:"Title"`
	ImageKey         string `json:"ImageKey"` // Use as unique ID if FileName is empty
	ArchivedMD5      string `json:"ArchivedMD5"`
	ArchivedSize     int64  `json:"ArchivedSize"`
	ArchivedUri      string `json:"ArchivedUri"`
	IsVideo          bool   `json:"IsVideo"`
	Processing       bool   `json:"Processing"`
	UploadKey        string `json:"UploadKey"`
	DateTimeOriginal string `json:"DateTimeOriginal"`
	DateTimeUploaded string `json:"DateTimeUploaded"`
	Caption          string `json:"Caption"`
	Keywords         string `json:"Keywords"`
	Latitude         string `json:"Latitude"`
	Longitude        string `json:"Longitude"`
	Altitude         string `json:"Altitude"`
	Format           string `json:"Format"`
	Width            int    `json:"Width"`
	Height           int    `json:"Height"`
	OriginalWidth    int    `json:"OriginalWidth"`
	OriginalHeight   int    `json:"OriginalHeight"`
	Size             int64  `json:"Size"`
	Hidden           bool   `json:"Hidden"`
	Watermarked      bool   `json:"Watermarked"`
	Collectable      bool   `json:"Collectable"`
	IsArchive        bool   `json:"IsArchive"`
	Status           string `json:"Status"`
	SubStatus        string `json:"SubStatus"`
	WebUri           string `json:"WebUri"`
	ThumbnailUrl     string `json:"ThumbnailUrl"`
	AlbumTitle       string // From album
	AlbumDescription string // From album
	AlbumKeywords    string // From album
	AlbumCreated     string // From album
	AlbumLastUpdated string // From album
	Uris             struct {
		ImageMetadata struct {
			Uri string `json:"Uri"`
		} `json:"ImageMetadata"`
		LargestVideo struct {
			Uri string `json:"Uri"`
		} `json:"LargestVideo"`
	} `json:"Uris"`

	builtFilename string // The final filename, after template replacements
}

func (a *albumImage) buildFilename(tmpl *template.Template) error {
	replacementVars := map[string]string{
		"FileName":      a.FileName,
		"ImageKey":      a.ImageKey,
		"ArchivedMD5":   a.ArchivedMD5,
		"UploadKey":     a.UploadKey,
		"Date":          "",
		"Time":          "",
		"Extension":     filepath.Ext(a.FileName),
		"FileNameNoExt": strings.TrimSuffix(filepath.Base(a.FileName), filepath.Ext(a.FileName)),
	}

	tm, err := time.Parse(time.RFC3339, a.DateTimeOriginal)
	if err == nil {
		replacementVars["Date"] = tm.Format(time.DateOnly)
		replacementVars["Time"] = strings.Replace(tm.Format(time.TimeOnly), ":", "_", -1)
	}

	var builtFilename bytes.Buffer
	if err := tmpl.Execute(&builtFilename, replacementVars); err != nil {
		return err
	}

	a.builtFilename = builtFilename.String()
	if a.builtFilename == "" {
		return errors.New("empty resulting name")
	}
	return nil
}

func (a *albumImage) Name() string {
	if a.builtFilename != "" {
		return a.builtFilename
	}

	if a.FileName != "" {
		return a.FileName
	}

	return a.ImageKey
}

type albumVideo struct {
	Response struct {
		LargestVideo struct {
			Size int64  `json:"Size"`
			Url  string `json:"Url"`
		} `json:"LargestVideo"`
	} `json:"Response"`
}
