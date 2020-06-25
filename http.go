package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

const maxRetries = 3
const baseAPIURL = "https://api.smugmug.com"

type header struct {
	name  string
	value string
}

func (s *smugmugHandler) get(url string, obj interface{}) error {
	return getJSON(fmt.Sprintf("%s%s", baseAPIURL, url), obj)
}

func getJSON(url string, obj interface{}) error {
	var result interface{}
	for i := 1; i <= maxRetries; i++ {
		resp, err := makeAPICall(url)
		if err != nil {
			return err
		}
		err = json.NewDecoder(resp.Body).Decode(&obj)
		defer resp.Body.Close()
		if err != nil {
			log.Errorf("%s: reading response\n[ERR] %s", url, err)
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

func addHeaders(req *http.Request, headers []header) {
	// Add headers if provided
	for _, h := range headers {
		req.Header.Add(h.name, h.value)
	}
}
