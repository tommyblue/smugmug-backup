package smugmug

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

// maxRetries defines the number of http calls retries (in case of errors) before giving up
const maxRetries = 3
const baseAPIURL = "https://api.smugmug.com"

type header struct {
	name  string
	value string
}

type handler struct{}

// get calls getJSON with the given url
func (s *handler) get(url string, obj interface{}) error {
	if url == "" {
		return errors.New("Can't get empty url")
	}
	return getJSON(fmt.Sprintf("%s%s", baseAPIURL, url), obj)
}

// download the resource (image or video) from the given url to the given destination, checking
// if a file with the same size exists (and skipping the download in that case)
func download(dest, downloadURL string, fileSize int64) error {
	if _, err := os.Stat(dest); err == nil {
		if sameFileSizes(dest, fileSize) {
			log.Debug("File exists with same size:", downloadURL)
			return nil
		}
	}
	log.Info("Getting ", downloadURL)

	response, err := makeAPICall(downloadURL)
	if err != nil {
		return fmt.Errorf("%s: download failed with: %s", downloadURL, err)
	}
	defer response.Body.Close()

	// Create empty destination file
	file, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("%s: file creation failed with: %s", dest, err)
	}
	defer file.Close()

	// Copy the content to the file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return fmt.Errorf("%s: file content copy failed with: %s", dest, err)
	}

	log.Info("Saved ", dest)
	return nil
}

// getJSON makes a http calls to the given url, trying to decode the JSON response on the given obj
func getJSON(url string, obj interface{}) error {
	var result interface{}
	for i := 1; i <= maxRetries; i++ {
		log.Debug("Calling ", url)
		resp, err := makeAPICall(url)
		if err != nil {
			return err
		}
		err = json.NewDecoder(resp.Body).Decode(&obj)
		defer resp.Body.Close()
		if err != nil {
			log.Errorf("%s: reading response. %s", url, err)
			if i >= maxRetries {
				return err
			}
		} else {
			obj = result
			break
		}
	}
	return nil
}

// makeAPICall performs an HTTP call to the given url, returning the response
func makeAPICall(url string) (*http.Response, error) {
	client := &http.Client{}

	var resp *http.Response
	var errorsList []error
	for i := 1; i <= maxRetries; i++ {
		req, err := http.NewRequest("GET", url, nil)

		// Auth header must be generate every time (nonce must change)
		h, err := authorizationHeader(url)
		if err != nil {
			panic(err)
		}
		headers := []header{
			{name: "Accept", value: "application/json"},
			{name: "Authorization", value: h},
		}
		log.Debug(headers)
		addHeaders(req, headers)

		r, err := client.Do(req)
		if err != nil {
			log.Debugf("#%d %s: %s\n", i, url, err)
			errorsList = append(errorsList, err)
			if i >= maxRetries {
				for _, e := range errorsList {
					log.Error(e)
				}
				return nil, errors.New("Too many errors")
			}
			// Go on and try again after a little pause
			time.Sleep(2 * time.Second)
			continue
		}

		if r.StatusCode >= 400 {
			errorsList = append(errorsList, errors.New(r.Status))
			if i >= maxRetries {
				for _, e := range errorsList {
					log.Error(e)
				}
				return nil, errors.New("Too many errors")
			}

			if r.StatusCode == 429 {
				// Header Retry-After tells the number of seconds until the end of the current window
				log.Error("Got 429 too many requests, let's try to wait 10 seconds...")
				log.Errorf("Retry-After header: %s\n", r.Header.Get("Retry-After"))
				time.Sleep(10 * time.Second)
			}
			continue
		}
		resp = r
		break

	}
	return resp, nil
}

// addHeaders to the provided http request
func addHeaders(req *http.Request, headers []header) {
	for _, h := range headers {
		req.Header.Add(h.name, h.value)
	}
}
