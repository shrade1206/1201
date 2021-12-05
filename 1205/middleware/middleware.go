package middleware

import (
	"net/http"
	"strings"
	"token/util"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 可將Token放在Header、Body、URI中，3選1
		// 預設前端回傳的Token放在Header裡的 "Authorization"，先判斷是否為空
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 0,
				"msg":  "參數無效",
			})
			c.Abort() //停止Middleware之後的func請求
			return
		}
		// 例 authHeader  = "Bearer Token"，parts[0] = "Bearer"，parts[1] = Token，把Bearer與Token分開來解析
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusOK, gin.H{
				"code": 0,
				"msg":  "authHeader格式有誤",
			})
			c.Abort()
			return
		}
		// 解析加密過的Token，取得設定好的資料
		Data, err := util.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "無效或過期的Token",
			})
			c.Abort()
			return
		}
		// 保存在c的上下文，配合Next()讓Middleware後續的func可以使用
		c.Set("username", Data.Username)
		c.Next()
	}
}
