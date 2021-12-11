package main

import (
	"log"
	"net/http"
	"token/controller"
	"token/middleware"
	redis_DB "token/redis"
	"token/session"

	"github.com/gin-contrib/sessions"

	"github.com/gin-gonic/gin"
)


func main() {
	err := redis_DB.InitRedis()
	if err != nil {
		log.Printf("InitRedis Error :%s", err.Error())
		return
	}
	defer redis_DB.Client.Close()
	err = session.Initsession()
	if err != nil {
		log.Printf("Initsession Error :%s", err.Error())
		return
	}
	r := gin.Default()
	r.Use(sessions.Sessions("mysession", session.Store))

	// store := cookie.NewStore([]byte("secret"))
	// r.Use(sessions.Sessions("mysession", store))
	// 測試------------------------------------
	checkUser := func(c *gin.Context) (string, bool) {
		var ok bool
		u, ok := c.Get("username")
		if !ok {
			return "", false
		}
		var username string
		username, ok = u.(string) // panic
		if !ok {
			return "", false
		}
		return username, true
	}
	handler := func(c *gin.Context) {
		// 因使用c.next()，只要沒有關掉伺服器，註冊完即可使用，如果有關掉的話會報錯 Key "username" does not exist
		var username, ok = checkUser(c)
		if !ok {
			c.JSON(200, gin.H{
				"code": 0,
				"msg":  "not login",
			})
			return
		}

		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "success, " + username,
		})
	}
	r.GET("/auth01", handler)
	r.GET("/auth02", handler)
	// 註冊帳號--------------------------------
	r.POST("/Register", controller.Register)
	// 登錄帳號--------------------------------
	r.POST("/login", controller.Login)
	// ---------------------------------------
	authRouter := r.Group("").Use(middleware.JWTAuthMiddleware())
	{
		authRouter.GET("/logout", controller.Logout)
	}
	// NoRoute--------------------------------
	r.NoRoute(gin.WrapH(http.FileServer(http.Dir("./public"))), controller.Noroute)
	// --------------------------------
	err = r.Run(":8080")
	if err != nil {
		log.Fatal("8080 err : ", err.Error())
		return
	}
}
