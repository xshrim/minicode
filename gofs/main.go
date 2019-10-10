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

// 使用go module
// export GOPROXY=https://goproxy.io   // 设置module代理
// go mod init m        // 初始化module或者从已有项目迁移(生成go.mod)
// go mod tidy          // 更新依赖
// go mod vendor        // 将所有依赖库复制到本地vendor目录
// go run -m=vendor main.go
// go build -mod=vendor // 利用本地vendor中的库构建或运行
// go list -u -m all    // 列出所有依赖库
// go mod edit -fmt     // 格式化go.mod

// 交叉编译:
// CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o gofs.exe main.go  // windows
// CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gofs main.go    // linux

// 解决alpine镜像问题, udp问题, 时区问题
// RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2 && apk add -U util-linux && apk add -U tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime  # 解决go语言程序无法在alpine执行的问题和syslog不支持udp的问题和时区问题

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
