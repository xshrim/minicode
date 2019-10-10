package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// 交叉编译:
//CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o gofs.exe main.go

const maxUploadSize = 4 * (2 << 30) // 4 * 1GB
const filePath = "./"

// upload file
//curl -X POST -F "path=test" -F "file=@/home/xshrim/a.js" http://127.0.0.1:8080/upload
//curl -X POST -F "file=@/home/xshrim/a.js" http://127.0.0.1:8080/upload/test/a.js
func upload(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		// crutime := time.Now().Unix()
		// h := md5.New()
		// io.WriteString(h, strconv.FormatInt(crutime, 10))
		// token := fmt.Sprintf("%x", h.Sum(nil))
		t, _ := template.ParseFiles("front.html")
		// t.Execute(w, token)
		t.Execute(w, nil)
		return
	}

	r.ParseMultipartForm(maxUploadSize)

	fpath := strings.TrimSpace(r.FormValue("path"))

	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println("Retrieving file error: ", err.Error())
		w.WriteHeader(http.StatusNoContent)
		fmt.Fprintf(w, "Failed: "+err.Error())
		return
	}
	defer file.Close()

	log.Println(fmt.Sprintf("Uploading file [filename: %+v, filesize: %+vB, httpheader: %+v", handler.Filename, handler.Size, handler.Header))

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("Retrieving file error: ", err.Error())
		w.WriteHeader(http.StatusNoContent)
		fmt.Fprintf(w, "Failed: "+err.Error())
		return
	}

	// tempFile, err := ioutil.TempFile(filePath, handler.Filename)
	if fpath == "" {
		fpath = strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/upload"), handler.Filename)
	}

	if err := ioutil.WriteFile(filepath.Join(filePath, fpath, handler.Filename), fileBytes, os.ModePerm); err != nil {
		log.Println("Creating file error: ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed: "+err.Error())
		return
	}

	log.Println("Receiving file successfully")

	fmt.Fprintf(w, "success")

}

func main() {
	port := "8080"

	dir, err := filepath.Abs("./")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.FileServer(http.Dir(filePath)))
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/upload/", upload)

	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	log.Println(fmt.Sprintf("starting file server at folder:<%s> address:<0.0.0.0:%s>", dir, port))

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}

}
