package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"token/controller"
	"token/middleware"
	redis_DB "token/redis"
	"token/session"

	"github.com/gin-contrib/sessions"

	"github.com/gin-gonic/gin"
)

// type RegisterData struct {
// 	Username string `json:"username" form:"username" binding:"required,min=6,max=12"`
// 	Password string `json:"password" form:"password" binding:"required,min=6,max=12"`
// }
var id string

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
	// 可使用 "github.com/gin-contrib/sessions/redis" 跟redis連接
	r := gin.Default()
	r.Use(sessions.Sessions("mysession", session.Store))

	// store := cookie.NewStore([]byte("secret"))
	// r.Use(sessions.Sessions("mysession", store))
	// 測試------------------------------------
	r.GET("/test", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("username", "jiyuu")
		err = session.Save()
		if err != nil {
			log.Println(err.Error())
			return
		}
		id = session.ID()
		redis_DB.Client.Set(id, "jiyuu", 0)
		a := session.Get("username")
		fmt.Println(a)
		fmt.Println(session.ID())
		c.JSON(200, gin.H{
			"ID":       id,
			"Username": a,
		})
	})
	r.GET("test1", func(c *gin.Context) {
		b := redis_DB.Client.Get(id)
		fmt.Println(b)
	})
	// 測試認證--------------------------------
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
	authRouter := r.Group("/api").Use(middleware.JWTAuthMiddleware())
	{
		authRouter.POST("/logout", func(c *gin.Context) {
			c.JSON(200, controller.API_Error{
				Msg: "Token ok",
			})
		})
	}
	// NoRoute--------------------------------
	r.NoRoute(gin.WrapH(http.FileServer(http.Dir("./public"))), func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method
		fmt.Println(path)
		fmt.Println(method)
		//檢查path的開頭使是否為"/"
		if strings.HasPrefix(path, "/") {
			fmt.Println("Route ok")
		}
	})
	// --------------------------------
	err = r.Run(":8080")
	if err != nil {
		log.Fatal("8080 err : ", err.Error())
		return
	}
}
