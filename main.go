package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	timeout = 10 * time.Second
)

func main() {
	port := flag.Int("p", 8080, "port to use")
	flag.Parse()

	registerHandles()
	http.DefaultClient.Timeout = timeout
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		fmt.Println("Had an err", err)
	}
}

func registerHandles() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/upload", upload)
	http.Handle("/files/", http.StripPrefix("/files", http.FileServer(http.Dir("./files"))))

}

func homePage(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}
	http.ServeFile(w, req, "index.html")
}

func upload(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		// crutime := time.Now().Unix()
		// h := md5.New()
		// io.WriteString(h, strconv.FormatInt(crutime, 10))
		// token := fmt.Sprintf("%x", h.Sum(nil))

		// t, _ := template.ParseFiles("upload.gtpl")
		// t.Execute(w, token)
	} else {
		req.ParseMultipartForm(32 << 20)
		file, handler, err := req.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		err = req.ParseForm()
		fmt.Println(*req)
		var newFileName string
		baseDir := "./files/"
		dir := baseDir
		if err != nil {
			fmt.Println("failed to read form:", err)
			newFileName = handler.Filename
		} else {
			newFileName = req.FormValue("newName")
			dir += req.FormValue("dir") + "/"
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				if err := os.MkdirAll(dir, os.ModePerm); err != nil {
					dir = baseDir
				}
			}
		}
		f, err := os.OpenFile(dir+newFileName, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println("failed to open file:", newFileName, "with error=", err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
}
