package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"token/controller"
	"token/middleware"
	"token/redis"
	"token/util"

	"github.com/gin-gonic/gin"
)

type RegisterData struct {
	Username string `json:"username" form:"username" binding:"required,min=6,max=12"`
	Password string `json:"password" form:"password" binding:"required,min=6,max=12"`
}

func main() {
	err := redis.InitRedis()
	if err != nil {
		log.Printf("InitRedis Error :%s", err.Error())
		return
	}
	defer redis.Client.Close()

	r := gin.Default()
	// 測試------------------------------------
	r.GET("ttt1", func(c *gin.Context) {
		data, _ := c.Cookie("token")
		Data, err := util.ParseToken(data)
		if err != nil {
			return
		} else {
			fmt.Println(Data)
			c.SetCookie("token", data, -1, "/", "127.0.0.1", false, true)
		}
	})
	r.POST("test1", func(c *gin.Context) {
		var reg RegisterData
		err := c.ShouldBindJSON(&reg)
		if err != nil {
			log.Printf("Register BindJson Error : %s", err.Error())
			c.JSON(http.StatusOK, err.Error())
			return
		}
		a, _ := redis.Client.Set(reg.Username, reg.Password, 0).Result()
		b := redis.Client.Get(reg.Username)
		fmt.Println(len(a))
		fmt.Println(b)
	})
	// 測試生成--------------------------------
	r.GET("/test", func(c *gin.Context) {
		token, _ := util.GenToken("test")
		fmt.Println("Bearer " + token)
	})
	// 測試解析--------------------------------
	r.GET("/parse", func(c *gin.Context) {
		var tk util.TokenData
		err := c.BindJSON(&tk)
		if err != nil {
			log.Println("解析失敗")
			return
		}
		a, _ := util.ParseToken(tk.Token)
		fmt.Println(a.Username, a.Issuer, a.ExpiresAt)
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
