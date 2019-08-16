package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)
// 交叉编译:
//CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o gofs.exe main.go

func main() {
	port := "8080"

	dir, err := filepath.Abs("./")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.FileServer(http.Dir("./")))

	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	log.Println(fmt.Sprintf("starting file server at folder:<%s> address:<0.0.0.0:%s>", dir, port))

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}

}
