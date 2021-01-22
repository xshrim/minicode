package cpaasapi

import (
	"crypto/tls"
	"net/http"
	"time"

	"alertsyslog/src/config"

	"github.com/xshrim/gol"
)

func ComplexGetHttp(url string) (int, string) {
	gol.Info("正在请求平台项目信息...")

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 8 * time.Second,
	}
	//读取Api数据

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	authKey := "Authorization"
	authValue := "Token " + config.APITOKEN
	//req.Header.Add("Cookie","sessionid=5edb1f18c5a0cb334b42b2383c899e01")
	request.Header.Add(authKey, authValue)

	response, err := client.Do(request)

	if err != nil || response == nil {
		gol.Error("请求错误！ url :" + url + ", Token :" + config.APITOKEN[:5])
		return 404, ""
	}
	defer response.Body.Close()

	return handleResponse(response)
}

func handleResponse(respon *http.Response) (int, string) {
	var bodyStr string

	buf := make([]byte, 1024*4)

	for {
		n, err := respon.Body.Read(buf)
		if n == 0 {
			gol.Error("http response body read err :", err)
			break
		}
		bodyStr += string(buf[:n])
	}
	return respon.StatusCode, bodyStr
}
