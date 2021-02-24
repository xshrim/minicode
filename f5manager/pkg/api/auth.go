package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xshrim/f5m/pkg/model"
	"github.com/xshrim/f5m/pkg/utils"
)

func GetAuth(c *gin.Context) {
	code := 200
	data := make(map[string]interface{})
	username := c.Query("username")
	password := c.Query("password")

	if model.CheckAuth(username, password) {
		token, err := utils.GenerateToken(username, password)
		if err != nil {
			code = 500
		} else {
			data["token"] = token
		}
	} else {
		code = 401
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"data": data,
	})
}
