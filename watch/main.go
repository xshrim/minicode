package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	KEEP   = "keep"
	CREATE = "create"
	MODIFY = "modify"
	DELETE = "delete"
)

var changed bool

func GetArgs() (string, string, string) {
	var fpath, opt, cmd string
	flag.StringVar(&fpath, "f", "", "filepath")
	flag.StringVar(&opt, "e", "", "option")

	flag.Parse()

	cmd = strings.Join(flag.Args(), " ")

	// file
	if fpath == "" {
		fname := os.Getenv("WATCHFILE")
		cwd, err := os.Getwd()
		if fname != "" && err == nil {
			if strings.HasPrefix(fname, cwd) {
				fpath = fname
			} else {
				fpath = filepath.Join(cwd, fname)
			}
		}
	}
	if fpath == "" {
		fpath = "./config.json"
	}

	// option
	if opt == "" {
		opt = os.Getenv("WATCHOPT")
	}
	if opt == "" {
		opt = "modify"
	}

	// command
	if cmd == "" {
		cmd = os.Getenv("WATCHCMD")
	}

	return fpath, opt, cmd
}

func FileExists(filename string) bool {
	fi, err := os.Stat(filename)
	return (err == nil || os.IsExist(err)) && !fi.IsDir()
}

func FileModTime(filename string) int64 {
	file, err := os.Stat(filename)
	if err != nil {
		return time.Now().Unix()
	}

	return file.ModTime().Unix()
}

func Watch() {
	var action string

	fpath, opt, cmd := GetArgs()

	existed := FileExists(fpath)
	modtime := FileModTime(fpath)

	log.Printf("watching action %s for file %s with command %s", opt, fpath, cmd)

	ticker := time.NewTicker(time.Duration(time.Second * 1))
	defer ticker.Stop()
	for {
		<-ticker.C
		if !FileExists(fpath) {
			if existed {
				action = DELETE
				modtime = FileModTime(fpath)
			} else {
				action = KEEP
			}
			existed = false
		} else {
			if !existed {
				action = CREATE
				existed = true
				modtime = FileModTime(fpath)
			} else {
				newtime := FileModTime(fpath)
				if modtime != newtime {
					action = MODIFY
					modtime = newtime
				} else {
					action = KEEP
				}
			}
		}

		if action != KEEP {
			changed = true
		}

		if strings.Contains(opt, action) && cmd != "" {
			log.Printf("action: <%s> | command: <%s>", action, cmd)
			go RunCmd(cmd)
		}
	}
}

func RunCmd(cmd string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	command := exec.Command("sh", "-c", cmd)
	command.Stdout = &stdout
	command.Stderr = &stderr

	if err := command.Run(); err != nil {
		log.Println(err)
	}

	fmt.Println(stdout.String(), stderr.String())
}

func main() {

	go Watch()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
		if changed {
			changed = false
			// http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			http.Error(w, "File Changed", http.StatusNotFound)
		} else {
			// fmt.Fprintf(w, http.StatusText(http.StatusOK))
			fmt.Fprintf(w, "File Not Changed")
		}
	})
	log.Println("Listening on localhost:2233")
	log.Fatal(http.ListenAndServe(":2233", nil))
}
