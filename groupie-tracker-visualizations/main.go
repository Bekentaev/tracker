package main

import (
	"log"
	"net/http"

	"art/home"
)

func main() {
	mux := http.ServeMux{}
	mux.HandleFunc("/", home.Artists)
	mux.HandleFunc("/artists/", home.ArtistsPage)
	fileServer := http.FileServer(http.Dir("./ui/css/"))
	mux.Handle("/ui/css/", http.StripPrefix("/ui/css", fileServer))
	log.Println("Starting server on:localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", &mux))
}
