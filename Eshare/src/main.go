package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

func GenerateSID() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func Index(r *gin.Engine, c *gin.Context) {
	sid, err := c.Cookie("sid")
	fmt.Println(sid)
	session := sessions.Default(c)
	svalue := session.Get(sid)
	if err != nil || sid == "" || svalue == nil {
		// c.Request.URL.Path = "/login"
		// r.HandleContext(c)
		c.Redirect(http.StatusMovedPermanently, "/login")
	} else {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":    "ok",
			"username": svalue.(string),
		})
	}
}

func Login(c *gin.Context) {
	sid, err := c.Cookie("sid") // 根据用户浏览器的cookie中的特定关键字得到cookie id
	session := sessions.Default(c)
	svalue := session.Get(sid)                    // 根据cookie id查找该用户在session中的信息
	if err == nil && sid != "" && svalue != nil { // 如果用户仍然在session中, 则无需登录
		c.Redirect(http.StatusMovedPermanently, "/index")
	} else {
		username := c.PostForm("user")
		password := c.PostForm("passwd")

		if username == "admin" && password == "admin" {
			sid := GenerateSID()                                         // 随机生成一个cookie id
			c.SetCookie("sid", sid, 3600, "/", "localhost", false, true) // 设置用户浏览器客户端的cookie
			session.Set(sid, username)                                   // 以cookie id为键, 用户名为值(可以是任何数据类型)在session中存储用户信息
			// session.Options(sessions.Options{   // seesion 可选项设置(过期时间等)
			// })
			session.Save() // 保存session

			// c.Redirect(http.StatusMovedPermanently, "http://localhost:8080/index")
			c.Redirect(http.StatusMovedPermanently, "/index")
		} else {
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{
				"title": "login",
			})
		}
	}
}

func main() {
	router := gin.Default()
	// 表示在html或go代码中访问/static这个路由下的文件就是访问文件系统./static目录下的文件, 相对路径是相对main.go所在目录
	// 如index.html中访问/static/refs/js/jquery.js即表示访问./static/refs/js/jquery.js
	router.Static("/static", "./static")
	// 表示加载./static/templates目录下的所有文件(不包括子目录), 该目录下的文件可无需路径直接访问
	// 如需包含子目录则需使用"./static/templates/**/*"
	router.LoadHTMLGlob("./static/templates/*")
	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")

	store := memstore.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	router.GET("/index", func(c *gin.Context) {
		Index(router, c)
	})
	// router.POST("/login", Login)
	router.Any("/login", Login)
	//router.POST("/login", Login)

	/*
		router.GET("/index", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"title": "Main website",
			})
		})
	*/
	router.Run(":8080")
}
