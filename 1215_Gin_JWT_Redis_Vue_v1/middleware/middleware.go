package middleware

import (
	"token/controller"
	"token/redis"
	"token/util"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 直接讀取前端cookie裡的token
		token, err := c.Cookie("token")
		if err != nil {
			controller.ErrMsg(c, controller.Code_Param_Invalid, "請重新登錄", nil, err)
			return
		}
		// 解析加密過的Token，取得SessionID
		Data, err := util.ParseToken(token)
		if err != nil {
			// 移除過期Token
			controller.ErrMsg(c, controller.Code_Param_Invalid, "Token Expired", nil, err)
			c.SetCookie("token", token, -1, "/", "localhost", false, true)
			if err != nil {
				controller.ErrMsg(c, controller.Code_DB_Conn, "Redis", nil, err)
				return
			}
			c.Abort()
			return
		}
		// 從Redis取出username
		username, err := redis.Client.Get(Data.SessionID).Result()
		if username == Data.Username {
			// 保存在c的上下文，配合Next()讓Middleware後續的func可以使用
			c.Set("sessionid", Data.SessionID)
			c.Set("username", Data.Username)
			c.Next()
		} else {
			controller.ErrMsg(c, controller.Code_Param_Invalid, "資料有誤，請重新登錄", nil, err)
			c.Abort()
			return
		}
	}
}
