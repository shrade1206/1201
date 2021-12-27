package middleware

import (
	"net/http"
	"strings"
	"token/redis"
	"token/util"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 直接讀取前端cookie裡的token
		token, err := c.Cookie("token")
		if err != nil {
			util.ErrMsg(c, util.Code_Param_Invalid, "請重新登錄", nil, err)
			return
		}
		// 解析加密過的Token，取得SessionID
		Data, err := util.ParseToken(token)
		if err != nil {
			// 如果錯誤是expired就移除
			expired := strings.Contains(err.Error(), "Token is expired")
			if expired {
				util.ErrMsg(c, util.Code_Param_Invalid, "Token Expired", nil, err)
				err = redis.Client.Del(Data.SessionID).Err()
				if err != nil {
					util.ErrMsg(c, util.Code_DB_Conn, "Redis", nil, err)
					return
				}
			}
			if !expired {
				util.ErrMsg(c, util.Code_Param_Invalid, "資料有誤，請重新登錄", nil, err)
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
			util.ErrMsg(c, util.Code_Param_Invalid, "資料有誤，請重新登錄", nil, err)
			return
		}
	}
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			//主要設定Access-Control-Allow-Origin
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Set("content-type", "application/json")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
