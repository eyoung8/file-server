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
	baseDir = "./files/"
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
		http.Redirect(w, req, "/", http.StatusSeeOther)
	} else {
		req.ParseMultipartForm(32 << 20)
		file, handler, err := req.FormFile("uploadfile")
		if err != nil {
			fmt.Println("Failed to open form file:", err)
			io.WriteString(w, "Failed to open file")
			return
		}
		defer file.Close()
		if handler.Filename == "" {
			return
		}
		if err = req.ParseForm(); err != nil {
			fmt.Println("failed to read form, request:", *req, "with error:", err)
			io.WriteString(w, "Failed to read form")
			return
		}
		newFileName := req.FormValue("newName")
		if newFileName == "" {
			newFileName = handler.Filename
		}
		dir := getFullPath(baseDir, req.FormValue("dir"))
		err = makeDirectoriesIfNecessary(dir)
		if err != nil {
			io.WriteString(w, "Error making requested directory")
			fmt.Println("Error making directory:", err)
			return
		}
		if err = copyFileToDisk(file, dir, newFileName); err != nil {
			fmt.Println("Error copying file to disk:", err)
			io.WriteString(w, "Error copying file to disk")
			return
		}
		http.ServeFile(w, req, "upload.html")
	}
}

func getFullPath(base string, dir string) string {
	return base + dir + "/"
}

func makeDirectoriesIfNecessary(dir string) error {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			fmt.Println("Couldn't make all dirs:", err)
			return err
		}
	} else if err != nil {
		fmt.Println("Couldn't stat dir:", err)
		return err
	}
	return nil
}

func copyFileToDisk(file io.Reader, dir string, newFileName string) error {
	f, err := os.OpenFile(dir+newFileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("failed to open file:", newFileName, "with error=", err)
		return err
	}
	defer f.Close()
	if _, err = io.Copy(f, file); err != nil {
		fmt.Println("failed to copy file:", err)
		return err
	}
	return nil
}
