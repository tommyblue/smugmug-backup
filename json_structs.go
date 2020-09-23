package smugmug

import (
	"bytes"
	"errors"
	"html/template"
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
	URLPath string `json:"UrlPath"`
	Uris    struct {
		AlbumImages struct {
			URI string `json:"Uri"`
		} `json:"AlbumImages"`
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

type albumImage struct {
	AlbumPath    string // From album.URLPath
	FileName     string `json:"FileName"`
	ImageKey     string `json:"ImageKey"` // Use as unique ID if FileName is empty
	ArchivedMD5  string `json:"ArchivedMD5"`
	ArchivedSize int64  `json:"ArchivedSize"`
	ArchivedUri  string `json:"ArchivedUri"`
	IsVideo      bool   `json:"IsVideo"`
	Processing   bool   `json:"Processing"`
	UploadKey    string `json:"UploadKey"`
	Uris         struct {
		LargestVideo struct {
			Uri string `json:"Uri"`
		} `json:"LargestVideo"`
	} `json:"Uris"`

	builtFilename string // The final filename, after template replacements
}

func (a *albumImage) buildFilename(tmpl *template.Template) error {
	replacementVars := map[string]string{
		"FileName":    a.FileName,
		"ImageKey":    a.ImageKey,
		"ArchivedMD5": a.ArchivedMD5,
		"UploadKey":   a.UploadKey,
	}

	var builtFilename bytes.Buffer
	if err := tmpl.Execute(&builtFilename, replacementVars); err != nil {
		return err
	}

	a.builtFilename = builtFilename.String()
	if a.builtFilename == "" {
		return errors.New("Empty resulting name")
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
