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
		// 先用Token取出Session
		sessionId, err := redis.Client.Get(token).Result()
		if err != nil {
			controller.ErrMsg(c, controller.Code_Param_Invalid, "Param", nil, err)
			return
		}
		// 確認SessionID是否一致
		if sessionId == controller.Session.ID() {
			// 確認UserName是否一致
			username := controller.Session.Get("username")
			if err != nil {
				controller.ErrMsg(c, controller.Code_Session_Invalid, "SessionId", nil, err)
				return
			}
			// 解析加密過的Token，取得設定好的資料
			Data, err := util.ParseToken(token)
			if err != nil {
				controller.ErrMsg(c, controller.Code_Param_Invalid, "Token Expired", nil, err)
				// 移除過期Token
				c.SetCookie("token", token, -1, "/", "localhost", false, true)
				err = redis.Client.Del(Data.Username + "Token").Err()
				if err != nil {
					controller.ErrMsg(c, controller.Code_DB_Conn, "Redis", nil, err)
					return
				}
				return
			}
			// 確認username
			if username == Data.Username {
				// 保存在c的上下文，配合Next()讓Middleware後續的func可以使用
				c.Set("username", Data.Username)
				c.Next()
			} else {
				controller.ErrMsg(c, controller.Code_Param_Invalid, "資料有誤，請重新登錄", nil, err)
				c.Abort()
				return
			}
		}
	}
}
