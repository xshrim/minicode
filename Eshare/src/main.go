package main

import (
	"context"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"./doc"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"

	// go get go.mongodb.org/mongo-driver/mongo
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	//"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	url    = "mongodb://localhost:27017"
	dbname = "eshare"
)

func ReplaceStr(str string) string {
	replacer := strings.NewReplacer("\\", "", " ", "", "\n", "", "\r", "", "\t", "", "/", "", ":", "", "*", "", "?", "", "|", "", "<", "", ">", "", "@", "")
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

func Conn(url string) *mongo.Client {
	if url == "" {
		url = "url"
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		log.Println(err)
		return nil
	}

	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Println(err)
		return nil
	}

	return client

}

// 保存对象到数据库
func Save(obj interface{}, colname string) string {

	client := Conn(url)
	if client == nil {
		return ""
	}
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	db := client.Database(dbname)
	col := db.Collection(colname)

	/*
		var val interface{}
		if colname == "document" {
			val = obj.(*doc.Document)
		} else if colname == "page" {
			val = obj.(*doc.Page)
		}
	*/
	res, err := col.InsertOne(ctx, obj)
	if err != nil {
		log.Println(err)
		return ""
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		return oid.Hex()
	}

	return ""
}

// 更新数据库上的对象
func Update(id, colname string, val bson.M) int64 {

	client := Conn(url)
	if client == nil {
		return 0
	}
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	db := client.Database(dbname)
	col := db.Collection(colname)

	/*
		var update bson.M
		json.Unmarshal([]byte(`{ "$set": {"pagenum": 2}}`), &update)
		res, err := col.UpdateOne(ctx, bson.M{"id": document.ID}, update)
	*/
	res, err := col.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": val})
	if err != nil {
		log.Println(err)
		return 0
	}

	return res.ModifiedCount
	/*
		if oid, ok := res.UpsertedID.(primitive.ObjectID); ok {
			return oid.Hex()
		}
	*/
}

// 查找数据库上的文档对象
func FindDocument(id string) *doc.Document {
	client := Conn(url)
	if client == nil {
		return nil
	}
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	//objID, _ := primitive.ObjectIDFromHex("5c7452c7aeb4c97e0cdb75bf")
	//value := collection.FindOne(ctx, bson.M{"_id": objID})

	filter := bson.M{"id": id}
	document := new(doc.Document)

	db := client.Database(dbname)
	col := db.Collection("document")
	err := col.FindOne(ctx, filter).Decode(document)
	if err != nil {
		return nil
	}
	log.Println(document)

	return document
}

// 查找数据库上的页面对象
func FindPage(id string, start, end int64) []*doc.Page {
	client := Conn(url)
	if client == nil {
		return nil
	}

	var pages []*doc.Page
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	//objID, _ := primitive.ObjectIDFromHex("5c7452c7aeb4c97e0cdb75bf")
	//value := collection.FindOne(ctx, bson.M{"_id": objID})

	filter := bson.M{"id": id, "number": bson.M{"$gte": start, "$lte": end}}

	db := client.Database(dbname)
	col := db.Collection("page")
	cursor, err := col.Find(ctx, filter)
	if err != nil {
		return nil
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		page := new(doc.Page)
		cursor.Decode(page)
		pages = append(pages, page)
	}

	return pages
}

// 文档格式转换
func Convert(document *doc.Document) bool {
	if document == nil {
		return false
	}
	isImg := false
	fullname := path.Join("static/files", document.Name)
	fullname, err := filepath.Abs(fullname)
	if err != nil {
		return false
	}
	ext := filepath.Ext(fullname)
	pdffile := fullname
	switch strings.ToLower(ext) {
	case ".doc", ".docx", ".ppt", ".pptx", ".xls", ".xlsx", ".wps", ".rtf", ".pps", ".ppsx", ".dps", ".odp", ".pot", ".et", ".ods":
		// office文档转pdf
		// 需要为服务器安装相应字体, 否则可能导致转换后页数与原文件不一致
		// 常用字体: simsun.ttf(宋体SimSun), simhei.ttf(黑体SimHei)
		err := doc.OfficeToPDF(fullname)
		log.Println(err)
		pdffile = strings.TrimSuffix(fullname, ext) + ".pdf"
	case ".epub", ".umd", ".chm", ".mobi", ".md", ".txt", ".azw3", ".fb2", ".htmlz", ".lit", ".lrf", ".pdb", ".pmiz", ".rb", ".snb", ".tcr", ".txtz":
		var err error
		pdffile, err = doc.FileToPDF(fullname)
		log.Println(err)
	case ".jpg", ".jpeg", ".bmp", ".gif", ".tiff", ".webp", ".png":
		isImg = true
		log.Println("image")
	}
	if PathExists(pdffile) {

		defer func() {
			// 删除临时生成的pdf文件
			if !strings.HasSuffix(pdffile, document.Name) {
				os.Remove(pdffile)
			}
		}()

		var pagenum int64
		if isImg {
			pagenum = 1
		} else {
			pagenum = doc.PageNumber(pdffile)
		}
		log.Println(pagenum)
		if res := Update(document.ID, "document", bson.M{"pagenum": pagenum}); res == 0 {
			return false
		}

		// 多线程处理
		for i := 1; i <= int(pagenum); i++ {
			var pngfile string
			if !isImg {
				/*
					svgfile := strings.TrimSuffix(fullname, document.Ext) + "-" + strconv.Itoa(i) + ".svg"
					osvgfile := pngfile + ".svg"
				*/
				pngfile = strings.TrimSuffix(fullname, ext) + "-" + strconv.Itoa(i) + ".png"
				// png图片生成
				doc.PDF2PNG(pdffile, pngfile, i, 256) // ghostscript自带图像质量选项
				// TODO 多线程转换? 通过pdfinfo命令获取pdf基本信息(包括页数)

				// png图片压缩
				/*
					pngfile, err := doc.CompressPNG(opngfile)
					log.Println(pngfile, err)
				*/

				// pngfile = opngfile // 不作压缩
			} else {
				pngfile, _ = doc.ImageToPNG(pdffile) //此时pcffile就是源文件路径, 直接对其进行格式转换
			}

			if !PathExists(pngfile) {
				log.Println("Document convert finished")
				continue // pdf已经拆解完
			}

			fileBytes, err := ioutil.ReadFile(pngfile)
			log.Println(err)

			// 删除临时生成的png文件
			if !strings.HasSuffix(pngfile, ext) {
				os.Remove(pngfile)
			}

			page := &doc.Page{
				ID:      document.ID,
				Name:    document.Name,
				Pagenum: pagenum,
				Prenum:  document.Prenum,
				Number:  int64(i),
				Content: fileBytes,
			}
			//document.PngPages = append(document.PngPages, filebytes)

			res := Save(page, "page")
			log.Println(res)
		}
	}
	return false
}

func Echo(c *gin.Context) {
	id := c.PostForm("id")

	c.JSON(http.StatusCreated, map[string]interface{}{
		"res": FindDocument(id),
	})
	//c.JSON(http.StatusOK, res)
}

func DocInfo(c *gin.Context) {
	id := c.PostForm("id")
	if id == "" {
		return
	}

	document := FindDocument(id)
	if document == nil {
		log.Println("Not Found Document")
	}
	c.JSON(http.StatusCreated, map[string]interface{}{
		"res": document,
	})
}

func PageInfo(c *gin.Context) {
	id := c.PostForm("id")
	number := c.PostForm("range")
	if id == "" || number == "" {
		return
	}
	num := strings.Split(number, "-")

	var start, end int64
	start, err := strconv.ParseInt(num[0], 10, 64)
	if err != nil {
		log.Println(err)
		return
	}
	if len(num) > 1 {
		end, err = strconv.ParseInt(num[1], 10, 64)
		if err != nil {
			log.Println(err)
			return
		}
	}

	var pages []*doc.Page
	for _, page := range FindPage(id, start, end) {
		if page.Number <= page.Prenum {
			pages = append(pages, page)
		}
	}
	if pages == nil && len(pages) == 0 {
		log.Println("Not Found Page")
	}
	c.JSON(http.StatusCreated, map[string]interface{}{
		"res": pages,
	})
}

func Share(c *gin.Context) {
	sid, err := c.Cookie("sid")
	log.Println("Share:" + sid)
	session := sessions.Default(c)
	svalue := session.Get(sid)
	if err != nil || sid == "" || svalue == nil {
		// c.Request.URL.Path = "/login"
		// r.HandleContext(c)
		c.Redirect(http.StatusFound, "/login?surl=share")
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
		c.Redirect(http.StatusFound, "/login?surl=upload")
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

		size := file.Size

		name = ReplaceStr(name)

		name = sid + "@" + name

		fullname := path.Join("static/files", name)

		ext := strings.ToLower(path.Ext(fullname))

		if PathExists(fullname) {
			c.String(http.StatusOK, fmt.Sprintf("'%s' uploade failed: filename is existed!", file.Filename))
			return
		}

		log.Println(name)

		err = c.SaveUploadedFile(file, fullname)
		if err != nil {
			c.String(http.StatusOK, fmt.Sprintf("'%s' uploade failed: '%s'!", file.Filename, err.Error()))
			return
		}

		if !PathExists(fullname) {
			c.String(http.StatusOK, fmt.Sprintf("'%s' uploade failed!", file.Filename))
			return
		}

		id := GetMD5(fullname)
		title := c.PostForm("title")
		catalog := c.PostForm("catalog")
		class := c.PostForm("class")
		subclass := c.PostForm("subclass")
		tagstr := c.PostForm("tag")
		desc := c.PostForm("desc")
		date := time.Now().Unix()

		var tags []string
		for _, tag := range strings.Split(tagstr, ",") {
			tags = append(tags, strings.TrimSpace(tag))
		}

		var perms [5]int64

		teamPermStr := strings.TrimSpace(c.PostForm("teamperm"))
		if teamPermStr == "" {
			perms[0] = int64(0)
		} else {
			perm, err := strconv.ParseInt(teamPermStr, 10, 64)
			if err != nil {
				perms[0] = int64(0)
			} else {
				perms[0] = perm
			}
		}

		orgPermStr := strings.TrimSpace(c.PostForm("orgperm"))
		if orgPermStr == "" {
			perms[1] = int64(0)
		} else {
			perm, err := strconv.ParseInt(orgPermStr, 10, 64)
			if err != nil {
				perms[1] = int64(0)
			} else {
				perms[1] = perm
			}
		}

		corpPermStr := strings.TrimSpace(c.PostForm("corpperm"))
		if corpPermStr == "" {
			perms[2] = int64(-1)
		} else {
			perm, err := strconv.ParseInt(corpPermStr, 10, 64)
			if err != nil {
				perms[2] = int64(-1)
			} else {
				perms[2] = perm
			}
		}

		grpPermStr := strings.TrimSpace(c.PostForm("grpperm"))
		if grpPermStr == "" {
			perms[3] = int64(-1)
		} else {
			perm, err := strconv.ParseInt(grpPermStr, 10, 64)
			if err != nil {
				perms[3] = int64(-1)
			} else {
				perms[3] = perm
			}
		}

		consPermStr := strings.TrimSpace(c.PostForm("consperm"))
		if consPermStr == "" {
			perms[4] = int64(-1)
		} else {
			perm, err := strconv.ParseInt(consPermStr, 10, 64)
			if err != nil {
				perms[4] = int64(-1)
			} else {
				perms[4] = perm
			}
		}

		prenum, err := strconv.ParseInt(c.PostForm("prenum"), 10, 64)
		if err != nil {
			prenum = -1
		}

		/*
			document := &doc.Document{
				ID:       id,
				Ext:      ext,
				Owner:    owner,
				Title:    title,
				Name:     name,
				Catalog:  catalog,
				Class:    class,
				SubClass: subclass,
				Tag:      tags,
				Desc:     desc,
				Size:     size,
				Pagenum:  0,
				Vcnt:     0,
				Dcnt:     0,
				Score:    30,
				Ccnt:     0,
				Rcnt:     0,
				Perm:     perms,
				Prenum:   prenum,
				Date:     date,
				Endorse:  "光大科技 运维服务部",
				Block:    300,
				Txhash:   "1adbf20ace0d1601b00cc2b9dfdd4a431cfff9a13f6a6f5e5e4a80c897e0f7a8",
				Status:   0,
			}
		*/

		document := &doc.Document{
			ID:       id,
			Ext:      ext,
			Owner:    owner,
			Title:    title,
			Name:     name,
			Catalog:  catalog,
			Class:    class,
			SubClass: subclass,
			Tag:      tags,
			Desc:     desc,
			Size:     size,
			Pagenum:  0,
			Vcnt:     12,
			Dcnt:     3,
			Score:    45,
			Ccnt:     1,
			Rcnt:     2,
			Perm:     perms,
			Prenum:   prenum,
			Date:     date,
			Endorse:  "光大科技 运维服务部",
			Block:    301,
			Txhash:   "fa5c6b15723ec5d0aa104cf943611ebdefeb0a201a25d99464806aaa8c9326d0",
			Status:   0,
		}

		log.Println(document)
		res := Save(document, "document")
		if res != "" {
			go Convert(document)
		}
		c.HTML(http.StatusOK, "upinfo.html", gin.H{
			"title":    "光大E享-成功",
			"username": svalue.(string),
			"docid":    document.ID,
			"blknum":   document.Block,
			"txid":     document.Txhash,
			"endorse":  document.Endorse,
		})
		//c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	}
}

func View(c *gin.Context) {
	//	docid := c.PostForm("docid")
	sid, err := c.Cookie("sid")
	log.Println("View:" + sid)
	session := sessions.Default(c)
	svalue := session.Get(sid)
	if err != nil || sid == "" || svalue == nil {
		// c.Request.URL.Path = "/login"
		// r.HandleContext(c)
		c.Redirect(http.StatusFound, "/login?surl=view")
	} else {
		c.HTML(http.StatusOK, "view.html", gin.H{
			"title":    "光大E享-浏览",
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
		c.Redirect(http.StatusFound, "/login?surl=index")
	} else {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":    "光大E享-主页",
			"username": svalue.(string),
		})
	}
}

func Login(c *gin.Context) {
	surl := c.Query("surl")
	log.Println(surl)
	sid, err := c.Cookie("sid") // 根据用户浏览器的cookie中的特定关键字得到cookie id
	log.Println("Login:" + sid)
	session := sessions.Default(c)
	svalue := session.Get(sid)
	// 根据cookie id查找该用户在session中的信息
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

			if surl != "" {
				c.Redirect(http.StatusFound, "/"+surl)
			} else {
				c.Redirect(http.StatusFound, "/index")
			}
		} else {
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{
				"title": "光大E享-登录",
				"surl":  surl,
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
		/*
			c.HTML(http.StatusUnauthorized, "test.html", gin.H{
				"title": "Test Page",
			})
		*/
		c.HTML(http.StatusOK, "upinfo.html", gin.H{
			"title":    "光大E享-成功",
			"username": "张三",
			"docid":    "2b9dfdd4a431cfff9a13f",
			"blknum":   300,
			"txid":     "1adbf20ace0d1601b00cc2b9dfdd4a431cfff9a13f6a6f5e5e4a80c897e0f7a8",
			"endorse":  "光大科技 运维服务部",
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

	router.POST("/echo", Echo)

	router.POST("/pageinfo", PageInfo)

	router.POST("/docinfo", DocInfo)

	router.GET("/", func(c *gin.Context) {
		Index(router, c)
	})

	/*
		router.GET("/index", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"title": "Main website",
			})
		})
	*/
	router.Run(":8080")
}
