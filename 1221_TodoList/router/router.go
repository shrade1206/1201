package router

import (
	"log"
	"net/http"
	"todoList/controller"
	"todoList/corsMiddleware"

	"github.com/gin-gonic/gin"
)

func Router() error {

	r := gin.Default()
	r.Use(corsMiddleware.Cors())
	// r.Use(cors.Default())
	cors := r.Group("").Use(corsMiddleware.Cors())
	{
		//---------------------------------------
		cors.POST("/insert", controller.Insert)
		//---------------------------------------
		cors.GET("/get", controller.Get)
		//---------------------------------------
		cors.GET("/getpage", controller.GetPage)
		//---------------------------------------
		cors.PUT("/put/:id", controller.Put)
		//---------------------------------------
		cors.DELETE("/del/:id", controller.Del)
		//---------------------------------------
	}
	r.NoRoute(gin.WrapH(http.FileServer(http.Dir("./public"))), controller.Router)

	err := r.Run(":8082")
	if err != nil {
		log.Fatal("8080 err : ", err.Error())
		return err
	}
	return nil
}
