package midware

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xshrim/f5m/pkg/global"
	"github.com/xshrim/f5m/pkg/utils"
	"github.com/xshrim/gol"
	"github.com/xshrim/gol/tk"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		code := global.SUCCESS
		var data interface{}

		token := c.Query("token")
		if token == "" {
			authHeader := c.Request.Header.Get("Authorization")
			if authHeader != "" {
				parts := strings.SplitN(authHeader, " ", 2)
				if !(len(parts) == 2 && parts[0] == "Bearer") {
					code = global.ERROR_AUTH_CHECK_TOKEN_FORMAT
				} else {
					token = parts[1]
				}
			} else {
				if jsonData, err := ioutil.ReadAll(c.Request.Body); err == nil {
					token = tk.Jsquery(string(jsonData), "token").(string)
				} else {
					code = global.ERROR_AUTH_TOKEN_NOT_EXIST
				}
			}
		}

		if token != "" {
			claims, err := utils.ParseTokenUnverified(token)
			if err != nil {
				code = global.ERROR_AUTH_CHECK_TOKEN_FORMAT
				gol.Error(global.Msgs[code], err)
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = global.ERROR_AUTH_CHECK_TOKEN_EXPIRE
				gol.Error(global.Msgs[code], err)
			} else {
				c.Set("id", claims.Kid)
				c.Set("email", claims.Email)
				// TODO check token in local database
			}
		}

		if code != global.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  global.Msgs[code],
				"data": data,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
