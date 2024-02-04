package main

import (
	_ "embed"
	"encoding/json"
	"net/http"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"
)

//go:embed authuser.json
var authuser []byte

//go:embed user.json
var user []byte

//go:embed useralbums_1.json
var useralbums_1 []byte

//go:embed albumimages_1.json
var albumimages_1 []byte

//go:embed photo.jpg
var photo []byte

func parseJson(content []byte) any {
	var dest interface{}
	if err := json.Unmarshal(content, &dest); err != nil {
		log.Fatal(err)
	}

	return dest
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	api := chi.NewRouter()

	api.Route("/v2", func(r chi.Router) {
		// User details. Get the first page of the user albums
		r.Get("/user/{username}", func(w http.ResponseWriter, r *http.Request) {
			//username := chi.URLParam(r, "username")
			//resp.Response.User.Uris.UserAlbums.URI = fmt.Sprintf("/api/v2/user/%s!albums", username)
			responseOk(w, parseJson(user))
		})

		r.Get("/album/{albumId}!images", func(w http.ResponseWriter, r *http.Request) {
			//albumId := chi.URLParam(r, "albumId")
			responseOk(w, parseJson(albumimages_1))
		})

		r.Get("/user/{username}!albums", func(w http.ResponseWriter, r *http.Request) {
			responseOk(w, parseJson(useralbums_1))
		})
	})

	api.Get("/v2!authuser", func(w http.ResponseWriter, r *http.Request) {
		responseOk(w, parseJson(authuser))
	})

	r.Mount("/api", api)

	r.Get("/photos/{fname}.jpg", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/jpg")
		w.Write(photo)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		log.Infof("404 URI => %s", r.RequestURI)
		w.WriteHeader(http.StatusNotFound)
	})

	addr := "127.0.0.1:3000"
	log.Infof("Starting server on %s", addr)
	http.ListenAndServe(addr, r)
}

func responseOk(w http.ResponseWriter, resp any) {
	w.Header().Set("Content-Type", "application/json")
	h := http.StatusOK
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Warnf("Error: %v", err)
		h = http.StatusInternalServerError
	}
	w.WriteHeader(h)
}
