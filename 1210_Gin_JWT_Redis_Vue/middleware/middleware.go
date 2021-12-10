package middleware

import (
	"log"
	"net/http"
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
			c.JSON(http.StatusOK, controller.API_Error{
				Code: controller.Code_Param_Invalid,
				Msg:  "Cookie Error :" + err.Error(),
			})
			c.Abort() //停止Middleware之後的func請求
			return
		}
		// 用Token取出session，確認username是否相同，ok才解Token
		sessionId, err := redis.Client.Get(token).Result()
		if err != nil {
			log.Printf("token")
			c.JSON(http.StatusOK, controller.API_Error{
				Code: controller.Code_Param_Invalid,
				Msg:  "Param Error" + err.Error(),
			})
			return
		}
		// 確認SessionID
		if sessionId == controller.Session.Id {
			username, err := redis.Client.Get(sessionId).Result()
			if err != nil {
				log.Printf("SessionId Error :%s", err.Error())
				c.JSON(http.StatusOK, controller.API_Error{
					Code: controller.Code_Param_Invalid,
					Msg:  "sessionId Error" + err.Error(),
				})
				return
			}
			// 確認username
			if username == controller.UserData.Username {
				// 解析加密過的Token，取得設定好的資料
				Data, err := util.ParseToken(token)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusOK, controller.API_Error{
						Code: controller.Code_Param_Invalid,
						Msg:  "Token Expired :" + err.Error(),
					})
					// 移除過期Token
					c.SetCookie("token", token, -1, "/", "localhost", false, true)
					redis.Client.Del(Data.Username + "Token")
					return
				}
				// 保存在c的上下文，配合Next()讓Middleware後續的func可以使用
				c.Set("username", Data.Username)
				c.Next()
			} else {
				log.Printf("UserName Error :%s", err.Error())
				c.JSON(http.StatusOK, controller.API_Error{
					Code: controller.Code_Param_Invalid,
					Msg:  "UserName Error" + err.Error(),
				})
				return
			}

		}
	}
}
