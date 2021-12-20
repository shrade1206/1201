package router

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"todoList/controller"

	"github.com/gin-gonic/gin"
)

func Router() error {
	r := gin.Default()
	//---------------------------------------
	r.POST("/insert", controller.Insert)
	//----------------------------
	r.GET("/get", controller.Get)
	//---------------------------------------
	r.GET("/getpage", controller.GetPage)
	//---------------------------------------
	r.PUT("/put/:id", controller.Put)
	//---------------------------------------
	r.DELETE("/del/:id", controller.Del)
	//---------------------------------------
	r.NoRoute(gin.WrapH(http.FileServer(http.Dir("./public"))), func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method
		fmt.Println(path)
		fmt.Println(method)
		//檢查path的開頭使是否為"/"
		if strings.HasPrefix(path, "/") {
			fmt.Println("ok")
		}
	})
	err := r.Run(":8080")
	if err != nil {
		log.Fatal("8080 err : ", err.Error())
		return err
	}
	return nil
}
