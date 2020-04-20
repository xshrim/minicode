package auths

import (
	"fmt"
	"ido/tools/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var Auths map[string]*OAuth

func init() {
	Auths = make(map[string]*OAuth)
}

// 获取token, 如果mode为'code'表示通过code获取访问token和刷新token, 如果mode为'refresh_token'则表示通过refresh_token获取新的访问token和刷新token
func getToken(mode, url, client, scope, code, redirect, grant, secret string) (string, string, string, int64, int64, error) {

	info := fmt.Sprintf("client_id=%s&scope=%s&%s=%s&redirect_uri=%s&grant_type=%s&client_secret=%s", client, scope, mode, code, redirect, grant, secret)

	resp, err := http.Post(url+"/token", "application/x-www-form-urlencoded", strings.NewReader(info))
	if err != nil {
		return "", "", "", 0, 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", 0, 0, err
	}

	jsonStr := string(body)

	ttype := json.GetString(jsonStr, "token_type")
	atoken := json.GetString(jsonStr, "access_token")
	rtoken := json.GetString(jsonStr, "refresh_token")
	texpire := json.GetInt(jsonStr, "expires_in")
	tcreate := time.Now().Unix()

	return ttype, atoken, rtoken, tcreate, texpire, nil
}

// TODO 需要定时遍历Auths, 调用此函数刷新Token
func UpdateToken(mode, state, str string) error {
	auth, ok := Auths[state]
	if !ok {
		return fmt.Errorf("state not exist")
	}

	gtype := "refresh_token"

	if mode == "code" { // 否则str为刷新token
		auth.Code = str
		gtype = "authorization_code"
	}

	ttype, atoken, rtoken, tcreate, texpire, err := getToken(mode, auth.Url, auth.Client, auth.Scope, str, auth.Redirect, gtype, auth.Secret)
	if err != nil {
		return err
	}

	auth.TokenType = ttype
	auth.AccessToken = atoken
	auth.RefreshToken = rtoken
	auth.TokenCreate = tcreate
	auth.TokenExpire = texpire

	fmt.Println("OAUTH INFO:", auth)
	return nil
}

func FetchToken(state string) (string, error) {
	auth, ok := Auths[state]
	if !ok {
		return "", fmt.Errorf("need access token")
	}

	current := time.Now().Unix()
	if auth.TokenCreate+auth.TokenExpire+10 > current {
		if err := UpdateToken("refresh_token", state, auth.RefreshToken); err != nil {
			return "", err
		}
	}

	return auth.AccessToken, nil
}
