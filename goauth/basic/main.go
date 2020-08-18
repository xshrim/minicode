package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

func checkAuth(w http.ResponseWriter, r *http.Request) bool {
	s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(s) != 2 {
		return false
	}

	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		return false
	}

	fmt.Println(string(b))

	pair := strings.SplitN(string(b), ":", 2)
	if len(pair) != 2 {
		return false
	}

	return pair[0] == "user" && pair[1] == "pass"
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("start")
	if checkAuth(w, r) {
		fmt.Println("Authorized")
		w.Write([]byte("hello world"))
		return
	}
	fmt.Println("unauthorized")

	w.Header().Set("WWW-Authenticate", "Basic realm='MY REALM'")
	w.WriteHeader(401)
	w.Write([]byte("401 Unauthorized\n"))
}

func main() {
	http.HandleFunc("/", index)

	http.ListenAndServe(":8888", nil)

}
