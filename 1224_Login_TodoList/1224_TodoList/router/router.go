package router

import (
	"net/http"
	"todoList/controller"
	"todoList/middleware"

	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
)

func Router() error {

	r := gin.Default()
	r.Use(middleware.Cors())
	authRouter := r.Group("").Use(middleware.JWTAuthMiddleware())
	{
		authRouter.POST("/insert", controller.Insert)
		//---------------------------------------
		authRouter.GET("/get", controller.Get)
		//---------------------------------------
		authRouter.GET("/getpage", controller.GetPage)
		//---------------------------------------
		authRouter.PUT("/put/:id", controller.Put)
		//---------------------------------------
		authRouter.DELETE("/del/:id", controller.Del)
	}
	r.NoRoute(gin.WrapH(http.FileServer(http.Dir("./public"))), controller.Router)

	err := r.Run(":8083")
	if err != nil {
		log.Fatal().Err(err).Msg("8083 Error")
		return err
	}
	return nil
}
