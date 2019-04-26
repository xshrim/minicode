package main

import (
	"context"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"./doc"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"

	// go get go.mongodb.org/mongo-driver/mongo
	"go.mongodb.org/mongo-driver/mongo"
	//"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func ReplaceStr(str string) string {
	replacer := strings.NewReplacer("\\", "", " ", "", "\n", "", "\t", "", "/", "", ":", "", "*", "", "?", "", "|", "", "<", "", ">", "", "@", "")
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
	return ReplaceStr(base64.URLEncoding.EncodeToString(b))
}

func GetMD5(filepath string) string {
	var filechunk uint64 = 8192 // we settle for 8KB

	file, err := os.Open(filepath)

	if err != nil {
		panic(err.Error())
	}

	defer file.Close()

	// calculate the file size
	info, _ := file.Stat()

	filesize := info.Size()

	blocks := uint64(math.Ceil(float64(filesize) / float64(filechunk)))

	hash := md5.New()

	for i := uint64(0); i < blocks; i++ {
		blocksize := int(math.Min(float64(filechunk), float64(filesize-int64(i*filechunk))))
		buf := make([]byte, blocksize)

		file.Read(buf)
		io.WriteString(hash, string(buf)) // append into the hash
	}

	return fmt.Sprintf("%x", hash.Sum(nil))
}

func Save(doc *doc.Document, url string) {

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Println(err)
	}

	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Println(err)
	}

	db := client.Database("eshare")
	col := db.Collection("document")
	res, err := col.InsertOne(ctx, doc)
	if err != nil {
		log.Println(err)
	}

	log.Println(res)

}

func Share(c *gin.Context) {
	log.Println("===================================")
	sid, err := c.Cookie("sid")
	log.Println("Share:" + sid)
	session := sessions.Default(c)
	svalue := session.Get(sid)
	if err != nil || sid == "" || svalue == nil {
		// c.Request.URL.Path = "/login"
		// r.HandleContext(c)
		c.Redirect(http.StatusFound, "/login")
	} else {
		c.HTML(http.StatusOK, "upload.html", gin.H{
			"title":    "光大E享-上传",
			"username": svalue.(string),
			"maxsize":  2 << 26,
			"reward":   5,
		})
	}
}

func Upload(c *gin.Context) {
	sid, err := c.Cookie("sid")
	log.Println("Upload:" + sid)
	session := sessions.Default(c)
	svalue := session.Get(sid)
	if err != nil || sid == "" || svalue == nil {
		// c.Request.URL.Path = "/login"
		// r.HandleContext(c)
		c.Redirect(http.StatusFound, "/login")
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

		name = sid + "@" + name

		path := path.Join("static/files", name)

		if PathExists(path) {
			c.String(http.StatusOK, fmt.Sprintf("'%s' uploade failed: filename is existed!", file.Filename))
			return
		}

		log.Println(name)

		err = c.SaveUploadedFile(file, path)
		if err != nil {
			c.String(http.StatusOK, fmt.Sprintf("'%s' uploade failed: '%s'!", file.Filename, err.Error()))
			return
		}

		if !PathExists(path) {
			c.String(http.StatusOK, fmt.Sprintf("'%s' uploade failed!", file.Filename))
			return
		}

		id := GetMD5(path)
		title := c.PostForm("title")
		catalog := c.PostForm("catalog")
		class := c.PostForm("class")
		subclass := c.PostForm("subclass")
		tag := c.PostForm("tag")
		desc := c.PostForm("desc")
		date := time.Now().Unix()

		price, err := strconv.ParseInt(c.PostForm("price"), 10, 64)
		if err != nil {
			price = 0
		}

		doc := &doc.Document{
			ID:       id,
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

		log.Println(doc)
		Save(doc, "mongodb://localhost:27017")
		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	}
}

func View(c *gin.Context) {
	//	docid := c.PostForm("docid")
	log.Println("===================================")
	sid, err := c.Cookie("sid")
	log.Println("View:" + sid)
	session := sessions.Default(c)
	svalue := session.Get(sid)
	if err != nil || sid == "" || svalue == nil {
		// c.Request.URL.Path = "/login"
		// r.HandleContext(c)
		c.Redirect(http.StatusFound, "/login")
	} else {
		c.HTML(http.StatusOK, "test.html", gin.H{
			"title":    "光大E享-上传",
			"username": svalue.(string),
			"maxsize":  2 << 26,
			"reward":   5,
		})
	}

}

func Index(r *gin.Engine, c *gin.Context) {
	sid, err := c.Cookie("sid")
	log.Println("Index:" + sid)
	session := sessions.Default(c)
	svalue := session.Get(sid)
	if err != nil || sid == "" || svalue == nil {
		// c.Request.URL.Path = "/login"
		// r.HandleContext(c)
		//c.Redirect(http.StatusMovedPermanently, "/login")
		c.Redirect(http.StatusFound, "/login")
	} else {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":    "光大E享-主页",
			"username": svalue.(string),
		})
	}
}

func Login(c *gin.Context) {
	sid, err := c.Cookie("sid") // 根据用户浏览器的cookie中的特定关键字得到cookie id
	log.Println("Login:" + sid)
	session := sessions.Default(c)
	svalue := session.Get(sid)                    // 根据cookie id查找该用户在session中的信息
	if err == nil && sid != "" && svalue != nil { // 如果用户仍然在session中, 则无需登录
		c.Redirect(http.StatusFound, "/index")
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
			// 用301状态码(StatusMovedPermanently)作重定向会导致路由失效(缓存问题), 改为302状态码正常(StatusFound)
			c.Redirect(http.StatusFound, "/index") 
		} else {
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{
				"title": "光大E享-登录",
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

	router.Any("/test", func(c *gin.Context) {
		c.HTML(http.StatusUnauthorized, "test.html", gin.H{
			"title": "Test Page",
		})
	})
	router.GET("/index", func(c *gin.Context) {
		Index(router, c)
	})
	// router.POST("/login", Login)
	router.Any("/login", Login)
	//router.POST("/login", Login)

	router.GET("/share", Share)

	router.POST("/upload", Upload)

	router.GET("/view", View)
	/*
		router.GET("/", func(c *gin.Context) {
			Index(router, c)
		})
	*/

	/*
		router.GET("/index", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"title": "Main website",
			})
		})
	*/
	router.Run(":8080")
}
