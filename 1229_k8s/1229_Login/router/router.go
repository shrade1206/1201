package router

import (
	"token/controller"
	"token/middleware"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func Router() error {
	r := gin.Default()
	r.Use(middleware.Cors())
	// 註冊帳號--------------------------------
	r.POST("/Register", controller.Register)
	// 登錄帳號--------------------------------
	r.POST("/login", controller.Login)
	// 需驗證才能使用---------------------------------------
	authRouter := r.Group("").Use(middleware.JWTAuthMiddleware())
	{
		authRouter.GET("/middlewareAuth")
		//---------------------------------------
		authRouter.GET("/logout", controller.Logout)
		//---------------------------------------
		authRouter.GET("/middleware", controller.LoginStruct)
	}

	err := r.Run(":8082")
	if err != nil {
		log.Fatal().Err(err).Msg("8082 Error")
		return err
	}
	return nil
}
