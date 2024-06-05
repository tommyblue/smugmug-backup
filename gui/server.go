package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		ret := `{"message": "Hello, World!"}`
		fmt.Fprintf(w, ret)
	})

	log.Info("Server started on port 8088")
	log.Fatal(http.ListenAndServe("0.0.0.0:8088", mux))
}
