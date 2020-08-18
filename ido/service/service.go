package service

import (
	"fmt"
	"ido/auths"
	"ido/consts"
	"ido/ea"
	"ido/requests"
	"ido/tools/gintool"
	"ido/tools/log"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var logger = log.GetLogger("ido.servie", log.INFO)

func index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func authcheck(c *gin.Context) {
	state := c.PostForm("user") + consts.STRCONNECTOR + c.PostForm("app")
	status := "no"
	if auth, ok := auths.Auths[state]; ok {
		if auth.AccessToken != "" && auth.RefreshToken != "" {
			status = "ok"
		}
	}
	// c.JSON(http.StatusOK, gin.H{"func": "authcheck", "resp": status})
	gintool.ResultOkMsg(c, status, "authcheck")
}

func authorize(c *gin.Context) {

	auth := &auths.OAuth{
		Url:      c.PostForm("url"),
		Client:   strings.TrimSpace(c.PostForm("client")),
		Redirect: consts.REDIRECT,
		Scope:    strings.TrimSpace(c.PostForm("scope")),
		State:    strings.TrimSpace(c.PostForm("user")) + consts.STRCONNECTOR + c.PostForm("app"),
		Secret:   strings.TrimSpace(c.PostForm("secret")),
	}

	auths.Auths[auth.State] = auth //以state标识用户

	resp := fmt.Sprintf("%s/authorize?client_id=%s&response_type=code&redirect_uri=%s&response_mode=query&scope=%s&state=%s", auth.Url, auth.Client, auth.Redirect, auth.Scope, auth.State)
	gintool.ResultOkMsg(c, resp, "authorize")
}

func redirect(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	c.Header("Content-Type", "text/html; charset=utf-8")

	if err := auths.UpdateToken("code", state, code); err != nil {
		logger.Error(err)
		gintool.ResultStr(c, `<span>`+err.Error()+`</span>`)
		return
	}
	gintool.ResultStr(c, `<script>window.close()</script>`) // 直接关闭认证完成后的重定向页面
}

func request(c *gin.Context) {
	user := strings.TrimSpace(c.PostForm("user"))
	app := c.PostForm("app")
	method := strings.ToUpper(c.PostForm("method"))
	url := strings.TrimSpace(c.PostForm("url"))
	body := strings.TrimSpace(c.PostForm("body"))

	state := user + consts.STRCONNECTOR + app

	token, err := auths.FetchToken(state)
	if err != nil {
		gintool.ResultFailData(c, err.Error(), "request")
		return
	}

	contentType := "application/x-www-form-urlencoded"
	if strings.HasPrefix(body, "{") || strings.HasPrefix(body, "[{") {
		contentType = "application/json"
	}

	data, err := requests.Request(method, url, contentType, token, body)
	if err != nil {
		gintool.ResultFailData(c, err.Error(), "request")
		return
	}

	gintool.ResultOkMsg(c, strings.TrimSpace(string(data)), "request")
}

func link(c *gin.Context) {
	// c.ShouldBindJson(&link)无法获取通过字符串方式传入request body的数据
	// 可以通过c.ShouldBindBodyWith(&link, binding.JSON)方式获取
	// 获取request header可以通过c.Request.Header得到map
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		gintool.ResultFailData(c, err.Error(), "link")
		return
	}

	link, err := ea.New(body)
	if err != nil {
		gintool.ResultFailData(c, err.Error(), "link")
		return
	}

	if link.Status == "run" {
		link.RunLink()
	}

	gintool.ResultOkMsg(c, "Link "+link.Id+" status: "+link.Status, "link")
}

func Server() {
	// gin.SetMode(gin.ReleaseMode)
	route := gin.Default()
	route.Static("/ido/static", "./static")
	// r.LoadHTMLGlob("web/**/*.html")
	route.LoadHTMLGlob("./template/*.html")

	// route.Use(log.Logger())
	route.GET("/ido", index)

	route.Any("/ido/redirect", redirect)

	route.POST("/ido/authorize", authorize)

	route.POST("/ido/authcheck", authcheck)

	route.POST("/ido/request", request)

	route.POST("/ido/link", link)

	if err := route.Run(":8888"); err != nil {
		logger.Panicf("Error running ido service %s\n", err)
	}
}
