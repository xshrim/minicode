package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// git克隆
// git clone -b 'v2.3.21' --single-branch --depth 1 <url>

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
// CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags netgo -o gofs main.go    // linux

// 解决alpine镜像问题, udp问题, 时区问题
// RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2 && apk add -U util-linux && apk add -U tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime  # 解决go语言程序无法在alpine执行的问题和syslog不支持udp的问题和时区问题

const maxUploadSize = 4 * (2 << 30) // 4 * 1GB
var dir, host, port string

const html = `
<!DOCTYPE html>
<html lang="en">

<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <meta http-equiv="X-UA-Compatible" content="ie=edge" />
  <title>File Uploader</title>
  <!-- <script src="./bfi.js"></script> -->
</head>

<body>
  <p><strong>CMD Method</strong></p>
  <p>curl -X POST -F "path=bar" -F "file=@/root/foo/sample.pdf" http://{{.Host}}:{{.Port}}/upload</p>
  <p><strong>WEB Method</strong></p>
  <form enctype="multipart/form-data" action="http://{{.Host}}:{{.Port}}/upload" method="post" target="iiframe">
    <input name="path" placeholder="(Optional) remote storage path" size="30" />
    <input type="file" name="file" size="30" />
    <input type="submit" value="Upload" />
    <label> ¦ </label>
    <a href="http://{{.Host}}:{{.Port}}"><button type="button">Browse</button></a>
  </form>
  <iframe id="iiframe" name="iiframe" frameborder="0" width="600px" height="50px" ></iframe>
  <!-- <iframe id="iiframe" name="iiframe" frameborder="0" style="display:none;"></iframe> -->
</body>

</html>
`

type Server struct {
	Host string
	Port string
}

func GetLocalIP() string {
	if addrs, err := net.InterfaceAddrs(); err == nil {
		for _, address := range addrs {
			// check the address type and if it is not a loopback the display it
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}

// upload file
//curl -X POST -F "path=test" -F "file=@/home/xshrim/a.js" http://127.0.0.1:8080/upload
//curl -X POST -F "file=@/home/xshrim/a.js" http://127.0.0.1:8080/upload/test/a.js
func upload(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		// crutime := time.Now().Unix()
		// h := md5.New()
		// io.WriteString(h, strconv.FormatInt(crutime, 10))
		// token := fmt.Sprintf("%x", h.Sum(nil))
		// t, _ := template.ParseFiles("front.html")

		t, _ := template.New("index").Parse(html)

		// t.Execute(w, token)
		t.Execute(w, &Server{
			Host: host,
			Port: port,
		})
		return
	}

	r.ParseMultipartForm(maxUploadSize)

	fpath := strings.TrimSpace(r.FormValue("path"))

	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println("Receive file error: ", err.Error())
		// w.WriteHeader(http.StatusNoContent)
		fmt.Fprintf(w, "✘ Failed: "+err.Error())
		return
	}
	defer file.Close()

	log.Println(fmt.Sprintf("Receiving file [filename: %+v, filesize: %+vB, httpheader: %+v", handler.Filename, handler.Size, handler.Header))

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("Receive file error: ", err.Error())
		w.WriteHeader(http.StatusNoContent)
		fmt.Fprintf(w, "✘ Failed: "+err.Error())
		return
	}

	// tempFile, err := ioutil.TempFile(filePath, handler.Filename)
	if fpath == "" {
		fpath = strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/upload"), handler.Filename)
	}

	// fmt.Println(dir, fpath, handler.Filename)
	fullpath := filepath.Join(dir, fpath, handler.Filename)

	os.MkdirAll(filepath.Dir(fullpath), os.ModePerm)

	if err := ioutil.WriteFile(fullpath, fileBytes, os.ModePerm); err != nil {
		log.Println("Create file error: ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "✘ Failed: "+err.Error())
		return
	}

	log.Println("Receive file successfully")

	fmt.Fprintf(w, "✔ Succeeded")

}

func main() {
	// var dport = flag.String("port", "2333", "server port")
	// var dpath = flag.String("dir", "./", "server path")
	flag.StringVar(&port, "p", "2333", "server port")
	flag.StringVar(&port, "port", "2333", "server port")
	flag.StringVar(&dir, "d", "./", "server path")
	flag.StringVar(&dir, "dir", "./", "server path")

	flag.Parse()

	dir, err := filepath.Abs(dir)
	if err != nil {
		log.Fatal(err)
	}

	host = GetLocalIP()

	http.Handle("/", http.FileServer(http.Dir(dir)))

	http.HandleFunc("/upload", upload)
	http.HandleFunc("/upload/", upload)

	log.Println(fmt.Sprintf("serve path: <%s>", dir))
	log.Println(fmt.Sprintf("browse url: <0.0.0.0:%s>[%s]", port, host))
	log.Println(fmt.Sprintf("upload url: <0.0.0.0:%s/upload>[%s]", port, host))
	// log.Println(fmt.Sprintf("starting file server at folder:<%s> address:<0.0.0.0:%s>", dir, port))

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}

}
