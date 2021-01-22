package wechatHandle

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"alertsyslog/src/config"

	"github.com/xshrim/gol"
)

type wechatEventPayload struct {
	Content string `json:"content"`
}

type wechatEvent struct {
	ToParty string             `json:"toparty"`
	ToUser  string             `json:"touser"`
	ToTag   string             `json:"totag"`
	AgentID string             `json:"agentid"`
	MsgType string             `json:"msgtype"`
	Text    wechatEventPayload `json:"text"`
}

type wechatToken struct {
	AccessToken string `json:"access_token"`
}

type wechatResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

var wechatConfig = config.WechatConfig

// var config = &wechatConfig{
// 	URL:      "http://qyapi.weixin.qq.com/cgi-bin/",
// 	ProxyURL: "http://25.2.20.18:8080",
// 	Crop:     "wl310ee8108c",
// 	AgentID:  "1000097",
// 	Secret:   "AuEJQWLx7qKrMuqb2UnSb2d_8NQ_RZEm-FDXRilwy2k",
// }

const contentTypeJSON = "application/json"

// 转发
func forwordWechat(content string) {
	req, err := http.NewRequest(http.MethodGet, wechatConfig.URL+"gettoken", nil)
	if err != nil {
		gol.Errorf("无法连接微信服务： %s", err)
		return
	}

	q := req.URL.Query()
	q.Add("corpid", wechatConfig.Crop)
	q.Add("corpsecret", wechatConfig.Secret)
	req.URL.RawQuery = q.Encode()

	client, err := newClientFromConfig(wechatConfig.ProxyURL)
	if err != nil {
		gol.Error(err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		gol.Errorf("无法获取微信token： %s", err)
		return
	}
	defer resp.Body.Close()

	requestBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		gol.Errorf("无法获取微信token： %s", err)
		return
	}

	var wechatToken wechatToken
	if err := json.Unmarshal(requestBytes, &wechatToken); err != nil {
		gol.Errorf("无法获取微信token： %s", err)
		return
	}

	if wechatToken.AccessToken == "" {
		gol.Errorf("微信配置中的组织CropId填写不正确： %s", wechatConfig.Crop)
		return
	}

	wc := &wechatEvent{
		AgentID: wechatConfig.AgentID,
		MsgType: "text",
		Text: wechatEventPayload{
			Content: content,
		},
	}

	if wechatConfig.ToTag != "" {
		wc.ToTag = wechatConfig.ToTag
	} else if wechatConfig.ToUser != "" {
		wc.ToUser = wechatConfig.ToUser
	} else {
		gol.Error("微信通知目标不能为空")
		return
	}

	// 将原本硬编码的部分改成可配置的
	url := wechatConfig.URL + "message/send?access_token=" + wechatToken.AccessToken

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(wc); err != nil {
		gol.Errorf("微信消息序列化失败： %s", buf.String())
		return
	}

	resp, err = post(client, url, contentTypeJSON, &buf)
	if err != nil {
		gol.Errorf("发送微信通知失败： %s", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		gol.Errorf("发送微信通知失败： HTTP状态码为 %d", resp.StatusCode)
		return
	}

	requestBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		gol.Errorf("发送微信通知失败： %s", err)
		return
	}

	var weResp wechatResponse
	if err := json.Unmarshal(requestBytes, &weResp); err != nil {
		gol.Errorf("发送微信通知失败： %s", err)
		return
	}

	if weResp.Code != 0 {
		gol.Errorf("发送微信通知失败： %s", weResp.Error)
		return
	}
}

func newClientFromConfig(ProxyURL string) (*http.Client, error) {
	client := http.Client{
		Timeout: time.Second * 10,
	}

	if ProxyURL != "" {
		proxyURL, err := url.Parse(ProxyURL)
		if err != nil {
			return nil, errors.New("代理地址填写有误：" + ProxyURL)
		}
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	}

	return &client, nil
}

func post(client *http.Client, url string, bodyType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", bodyType)
	return client.Do(req)
}
