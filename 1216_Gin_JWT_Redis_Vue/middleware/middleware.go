package middleware

import (
	"strings"
	"token/controller"
	"token/redis"
	"token/util"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
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
			// 如果錯誤是expired就移除
			expired := strings.Contains(err.Error(), "token is expired")
			if expired {
				log.Warn().Caller().Msg("Token")
				controller.ErrMsg(c, controller.Code_Param_Invalid, "Token Expired", nil, err)
				err = redis.Client.Del(Data.SessionID).Err()
				if err != nil {
					controller.ErrMsg(c, controller.Code_DB_Conn, "Redis", nil, err)
					return
				}
			}
			if !expired {
				controller.ErrMsg(c, controller.Code_Param_Invalid, "資料有誤，請重新登錄", nil, err)
			}
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
			return
		}
	}
}
