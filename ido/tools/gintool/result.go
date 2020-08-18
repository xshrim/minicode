package gintool

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	Succ = iota
	Fail
)

func ResultStr(c *gin.Context, str string) {
	c.String(http.StatusOK, str)
}

func ResultMap(c *gin.Context, m map[string]interface{}) {
	c.JSON(http.StatusOK, m)
}

func ResultMsg(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, gin.H{"code": Succ, "data": nil, "msg": msg})
}
func ResultOk(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"code": Succ, "data": data, "msg": ""})
}
func ResultOkMsg(c *gin.Context, data interface{}, msg string) {
	c.JSON(http.StatusOK, gin.H{"code": Succ, "data": data, "msg": msg})
}

func ResultFail(c *gin.Context, err interface{}) {
	c.JSON(http.StatusOK, gin.H{"code": Fail, "data": nil, "msg": err})
}

func ResultFailData(c *gin.Context, data interface{}, err interface{}) {
	c.JSON(http.StatusOK, gin.H{"code": Fail, "data": data, "msg": err})
}

type RespData struct {
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
	Code int         `json:"code"`
}

/**
  @param  total  :总数
  @param  pageNum : 当前页
*/
type RespDataList struct {
	RespData
	Total int64 `json:"total"`
}

func ResultList(c *gin.Context, data interface{}, total int64) {
	c.JSON(http.StatusOK, gin.H{"code": Succ, "data": data, "msg": "", "total": total})
}

type RespDataPager struct {
	RespData
	Pager interface{} `json:"pager"`
}

func ResultPageList(c *gin.Context, data interface{}, pager interface{}) {
	c.JSON(http.StatusOK, gin.H{"code": Succ, "data": data, "msg": "", "pager": pager})
}
