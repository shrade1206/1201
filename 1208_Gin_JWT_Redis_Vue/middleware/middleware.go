package middleware

import (
	"net/http"
	"token/redis"
	"token/util"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 直接讀取前端cookie裡的token
		token, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusOK, redis.API_Error{
				Code: redis.Code_Param_Invalid,
				Msg:  "Cookie Error :" + err.Error(),
			})
			c.Abort() //停止Middleware之後的func請求
			return
		}
		// 解析加密過的Token，取得設定好的資料
		Data, err := util.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, redis.API_Error{
				Code: redis.Code_Param_Invalid,
				Msg:  "Token Expired :" + err.Error(),
			})
			// 移除過期Token
			c.SetCookie("token", token, -1, "/", "127.0.0.1", false, true)
			redis.Client.Del(Data.Username + "Token")
			return
		}
		// 保存在c的上下文，配合Next()讓Middleware後續的func可以使用
		c.Set("username", Data.Username)
		c.Next()
	}
}
