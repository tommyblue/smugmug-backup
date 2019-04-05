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
	// Code    int    `json:"Code"`
	// Message string `json:"Message"`
}

type albumsResponse struct {
	Response struct {
		URI   string  `json:"Uri"`
		Album []album `json:"Album"`
		Pages struct {
			// Total          int    `json:"Total"`
			// Start          int    `json:"Start"`
			// Count          int    `json:"Count"`
			// RequestedCount int    `json:"RequestedCount"`
			// FirstPage      string `json:"FirstPage"`
			// LastPage string `json:"LastPage"`
			NextPage string `json:"NextPage"`
			// PrevPage       string `json:"PrevPage"`
		} `json:"Pages"`
	} `json:"Response"`
	// Code    int    `json:"Code"`
	// Message string `json:"Message"`
}

type album struct {
	URLPath string `json:"UrlPath"`
	Uris    struct {
		AlbumImages struct {
			URI string `json:"Uri"`
		} `json:"AlbumImages"`
		// AlbumDownload struct {
		// 	URI string `json:"Uri"`
		// } `json:"AlbumDownload"`
	} `json:"Uris"`
}

type albumImagesResponse struct {
	Response struct {
		URI        string       `json:"Uri"`
		AlbumImage []albumImage `json:"AlbumImage"`
		Pages      struct {
			// Total          int    `json:"Total"`
			// Start          int    `json:"Start"`
			// Count          int    `json:"Count"`
			// RequestedCount int    `json:"RequestedCount"`
			// FirstPage      string `json:"FirstPage"`
			// LastPage string `json:"LastPage"`
			NextPage string `json:"NextPage"`
			// PrevPage       string `json:"PrevPage"`
		} `json:"Pages"`
	} `json:"Response"`
	// Code    int    `json:"Code"`
	// Message string `json:"Message"`
}

type albumImage struct {
	FileName    string `json:"FileName"`
	ArchivedUri string `json:"ArchivedUri"`
	IsVideo     bool   `json:"IsVideo"`
	Uris        struct {
		LargestVideo struct {
			Uri string `json:"Uri"`
		} `json: "Uris"`
		// 	LargestImage struct {
		// 		URI string `json:"Uri"`
		// 	} `json:"LargestImage"`
		// 	// AlbumDownload struct {
		// 	// 	URI string `json:"Uri"`
		// 	// } `json:"AlbumDownload"`
	} `json:"Uris"`
}

type albumVideo struct {
	Response struct {
		LargestVideo struct {
			Url string `json:"Url"`
		} `json:"LargestVideo"`
	} `json:"Response"`
}
