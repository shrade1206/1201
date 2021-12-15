package main

import (
	"net/http"
	"os"
	"token/controller"
	"token/middleware"
	redis_DB "token/redis"
	"token/session"
	"token/util"

	"github.com/gin-contrib/sessions"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

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
			return user, err
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
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
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
		err = redis_DB.Client.Del("1").Err()
		if err != nil {
			log.Warn().Caller().Err(err).Msg("GG")
		}

		_, err := util.ParseToken("1231")
		if err != nil {
			log.Warn().Caller().Err(err)
		}

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
		// 因使用c.next()，set完要的資料，通過middleware即可使用，如果有關掉的話會報錯 Key "username" does not exist
		var username, ok = checkUser(c)
		if !ok {
			c.JSON(200, gin.H{
				"Code": 0,
				"Msg":  "not login",
			})
			return
		}

		c.JSON(200, gin.H{
			"Code": 1,
			"Msg":  username,
		})
	}
	r.GET("/auth01", middleware.JWTAuthMiddleware(), handler)
	r.GET("/middleware", middleware.JWTAuthMiddleware(), controller.LoginStruct)
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
		log.Fatal().Caller().Err(err).Str("8080", "Err")
		return
	}
}
