package api

import (
	"github.com/gin-gonic/gin"
	"github.com/xshrim/f5m/pkg/global"
)

func resp(c *gin.Context, httpCode, code int, err error, data interface{}) {
	msg := global.Msgs[code]
	if err != nil {
		msg += ": " + err.Error()
	}

	c.JSON(httpCode, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}
