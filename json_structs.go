package main

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
	FileName     string `json:"FileName"`
	ImageKey     string `json:"ImageKey"` // Use as unique ID if FileName is empty
	ArchivedSize int64  `json:"ArchivedSize"`
	ArchivedUri  string `json:"ArchivedUri"`
	IsVideo      bool   `json:"IsVideo"`
	Processing   bool   `json:"Processing"`
	Uris         struct {
		LargestVideo struct {
			Uri string `json:"Uri"`
		} `json: "Uris"`
	} `json:"Uris"`
}

func (a *albumImage) Name() string {
	name := a.FileName
	if name == "" {
		name = a.ImageKey
	}
	return name
}

type albumVideo struct {
	Response struct {
		LargestVideo struct {
			Size int64  `json:"Size"`
			Url  string `json:"Url"`
		} `json:"LargestVideo"`
	} `json:"Response"`
}
