package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/", homePage)
	http.Handle("/files", http.StripPrefix("/files", http.FileServer(http.Dir("./files"))))
	http.ListenAndServe(":8080", nil)
}

func homePage(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "index.html")
}
