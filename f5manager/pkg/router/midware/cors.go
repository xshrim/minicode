package midware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORS(c *gin.Context) {
	// 后续可以考虑从环境变量中读取允许的源
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Access-Control-Allow-Methods", "*")

	// option 请求统统200
	if c.Request.Method == http.MethodOptions {
		c.AbortWithStatus(200)
	}
}
