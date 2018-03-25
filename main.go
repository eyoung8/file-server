package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"
)

const (
	timeout = 10 * time.Second
)

func main() {
	port := flag.Int("p", 8080, "port to use")
	flag.Parse()

	http.DefaultClient.Timeout = timeout
	http.HandleFunc("/", homePage)
	http.Handle("/files/", http.StripPrefix("/files", http.FileServer(http.Dir("./files"))))
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}

func homePage(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}
	http.ServeFile(w, req, "index.html")
}
