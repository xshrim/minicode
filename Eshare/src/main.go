package main

import (
	"crypto/rand"
	"doc"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

func ReplaceStr(str string) string {
	replacer := strings.NewReplacer("\\", "", " ", "", "\n", "", "\t", "", "/", "", ":", "", "*", "", "?", "", "|", "", "<", "", ">", "")
	return replacer.Replace(str)
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func GenerateSID() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func Share(c *gin.Context) {
	sid, err := c.Cookie("sid")
	fmt.Println(sid)
	session := sessions.Default(c)
	svalue := session.Get(sid)
	if err != nil || sid == "" || svalue == nil {
		// c.Request.URL.Path = "/login"
		// r.HandleContext(c)
		c.Redirect(http.StatusMovedPermanently, "/login")
	} else {
		c.HTML(http.StatusOK, "upload.html", gin.H{
			"title":    "ok",
			"username": svalue.(string),
			"maxsize":  2 << 26,
		})
	}
}

func Upload(c *gin.Context) {
	sid, err := c.Cookie("sid")
	fmt.Println(sid)
	session := sessions.Default(c)
	svalue := session.Get(sid)
	if err != nil || sid == "" || svalue == nil {
		// c.Request.URL.Path = "/login"
		// r.HandleContext(c)
		c.Redirect(http.StatusMovedPermanently, "/login")
	} else {

		owner := svalue.(string)

		file, err := c.FormFile("file") // 单文档上传
		if err != nil {
			c.String(http.StatusOK, fmt.Sprintf("'%s' uploade failed: '%s'!", file.Filename, err.Error()))
		}
		/*
			form, _ := c.MultipartForm()  // 多文档上传
			files := form.File["upload[]"]
			file := files[0]
		*/
		name := file.Filename

		name = ReplaceStr(name)

		path := path.Join("static/files", name)

		if PathExists(path) {
			c.String(http.StatusOK, fmt.Sprintf("'%s' uploade failed: filename is existed!", file.Filename))
		}

		log.Println(name)

		err = c.SaveUploadedFile(file, path)
		if err != nil {
			c.String(http.StatusOK, fmt.Sprintf("'%s' uploade failed: '%s'!", file.Filename, err.Error()))
		}
		title := c.PostForm("title")
		catalog := c.PostForm("catalog")
		class := c.PostForm("class")
		subclass := c.PostForm("subclass")
		price := c.PostForm("price")
		tag := c.PostForm("tag")
		desc := c.PostForm("desc")
		date := strconv.FormatInt(time.Now().Unix(), 10)

		doc := &doc.Document{
			ID:       "",
			SliceID:  "",
			Owner:    owner,
			Title:    title,
			Name:     name,
			Catalog:  catalog,
			Class:    class,
			SubClass: subclass,
			Tag:      tag,
			Desc:     desc,
			PageSize: 0,
			Vcnt:     0,
			Dcnt:     0,
			Score:    3,
			RaterNum: 0,
			Price:    price,
			Date:     date,
			Status:   0,
		}
		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	}
}

func View(c *gin.Context) {
	//	docid := c.PostForm("docid")

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
			c.SetCookie("sid", sid, 3600, "/", "localhost", false, true) // 设置用户浏览器客户端的cookie, domain部分必须与网页url一致
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

	router.MaxMultipartMemory = 8 << 26 //限制文件大小512MB

	router.GET("/", func(c *gin.Context) {
		Index(router, c)
	})
	router.GET("/index", func(c *gin.Context) {
		Index(router, c)
	})
	// router.POST("/login", Login)
	router.Any("/login", Login)
	//router.POST("/login", Login)

	router.Any("/share", Share)

	router.POST("/upload", Upload)

	router.Any("/view", View)

	/*
		router.GET("/index", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"title": "Main website",
			})
		})
	*/
	router.Run(":8080")
}
