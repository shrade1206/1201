package main

import (
	"fmt"
	"log"
	"net/http"
	"token/controller"
	"token/middleware"
	redis_DB "token/redis"
	"token/session"

	"github.com/gin-contrib/sessions"

	"github.com/gin-gonic/gin"
)

func GetUsername(c *gin.Context, username string) (controller.User, error) {
	var user controller.User
	// Exists查詢key是否存在，回傳true 或是false
	sel, err := controller.RedisExists(c, "username")
	if err != nil {
		log.Printf("Error : %s", err.Error())
	}
	if sel == 1 {
		value, err := redis_DB.Client.Get(username).Result()
		if err != nil {
			log.Printf("select Error :%s", err.Error())
			return user, nil
		} else {
			user = controller.User{Username: username, Password: value}
			return user, nil
		}
	} else {
		log.Printf("查無此帳號")
	}
	return user, nil
}

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

	r.GET("/test", func(c *gin.Context) {

		sel, err := controller.RedisExists(c, "username")
		if err != nil {
			return
		}
		n, err := GetUsername(c, "abc")
		fmt.Println(n.Password)

		fmt.Println(sel, err)
		a, _ := GetUsername(c, "jiyuusama")
		fmt.Println(a.Username)
		b, _ := GetUsername(c, "t")
		fmt.Println(b.Username)
	})

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
