package main

import (
	"net/http"
	"time"
)

const (
	timeout = 10 * time.Second
)

func main() {
	http.DefaultClient.Timeout = timeout
	http.HandleFunc("/", homePage)
	http.Handle("/files/", http.StripPrefix("/files", http.FileServer(http.Dir("./files"))))
	http.ListenAndServe(":8080", nil)
}

func homePage(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "index.html")
}
