package main

import (
	"net/http"
	"path/filepath"
)

func mainPage(w http.ResponseWriter, r *http.Request) {
	filepath := filepath.Join("client-side", "html", "index.html")
	http.ServeFile(w, r, filepath)
}

func newEventPage(w http.ResponseWriter, r *http.Request) {
	filepath := filepath.Join("client-side", "html", "createNewEvent.html")
	http.ServeFile(w, r, filepath)
}

func main() {
	http.HandleFunc("/", mainPage)
	http.HandleFunc("/createNewEvent", newEventPage)
	http.ListenAndServe(":8080", nil)
}
