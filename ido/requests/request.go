package requests

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func Request(method, url, contentType, token, pl string) ([]byte, error) {
	var data []byte
	client := &http.Client{}

	var body io.Reader

	if contentType == "application/json" {
		body = bytes.NewBuffer([]byte(pl))
	} else if contentType == "application/x-www-form-urlencoded" {
		body = strings.NewReader(pl)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Accept-Language", "zh-cn,zh;q=0.8,en-us;q=0.5,en;q=0.3")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:12.0) Gecko/20100101 Firefox/12.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if strings.HasPrefix(strconv.Itoa(resp.StatusCode), "20") {
		switch resp.Header.Get("Content-Encoding") {
		case "gzip":
			reader, _ := gzip.NewReader(resp.Body)
			for {
				buf := make([]byte, 1024)
				n, err := reader.Read(buf)

				if err != nil && err != io.EOF {
					return nil, err
				}

				if n == 0 {
					break
				}
				data = append(data, buf...)
			}
		default:
			data, _ = ioutil.ReadAll(resp.Body)
		}
		return data, nil
	}

	return nil, fmt.Errorf("request failed with code: " + strconv.Itoa(resp.StatusCode))
}
