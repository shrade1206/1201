package router

import (
	"todoList/controller"
	"todoList/middleware"

	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
)

func Router() error {

	r := gin.Default()
	r.Use(middleware.Cors())
	r.GET("/asd", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"Msg": "ok2",
		})
	})
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
		//
	}

	err := r.Run(":8083")
	if err != nil {
		log.Fatal().Err(err).Msg("8083 Error")
		return err
	}
	return nil
}
