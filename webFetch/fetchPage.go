package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/proxy"
)

func DirectGet(iurl string) {

	response, err := http.Get(iurl)

	if err != nil {
		log.Println("response error:", err)
	}

	stdout := os.Stdout
	_, err = io.Copy(stdout, response.Body)

	if err != nil {
		log.Println("copy error:", err)
	}

}

func AgentGet(iurl string) {

	client := &http.Client{}

	request, err := http.NewRequest("GET", iurl, nil)

	if err != nil {
		log.Println("request error:", err)
	}

	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:23.0) Gecko/20100101 Firefox/23.0")

	response, err := client.Do(request)

	if err != nil {
		log.Println("response error:", err)
	}

	stdout := os.Stdout
	_, err = io.Copy(stdout, response.Body)

	if err != nil {
		log.Println("copy error:", err)
	}

}

func ProxyGet(iurl, iproxy, agent string) {

	var transport *http.Transport
	client := &http.Client{}

	arr := strings.Split(iproxy, ":")
	switch strings.ToLower(arr[0]) {
	case "http", "https":
		pro, _ := url.Parse(iproxy)
		transport = &http.Transport{
			Proxy: http.ProxyURL(pro),
		}
	case "socks", "socks4", "socks5":
		pro := arr[1][2:] + ":" + arr[2]
		dialer, _ := proxy.SOCKS5("tcp", pro, nil, proxy.Direct)
		transport = &http.Transport{
			Dial: dialer.Dial,
		}
	default:
		transport = nil
	}

	if transport != nil {
		client.Transport = transport
	}

	// response, err := client.Get(url)

	request, err := http.NewRequest("GET", iurl, nil)

	if err != nil {
		log.Println("request error:", err)
	}

	if agent != "" {
		if strings.ToLower(agent) == "default" {
			request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:23.0) Gecko/20100101 Firefox/23.0")
		} else {
			request.Header.Add("User-Agent", agent)
		}
	}

	response, err := client.Do(request)

	if err != nil {
		log.Println("response error:", err)
	}

	stdout := os.Stdout
	_, err = io.Copy(stdout, response.Body)

	if err != nil {
		log.Println("copy error:", err)
	}

}

func main() {

	iurl := "http://www.baidu.com"
	iurl = "https://www.zhongzilou.com/list/ming/1"

	// DirectGet(url)
	iproxy := "socks://127.0.0.1:12345"
	ProxyGet(iurl, iproxy, "default")
}
