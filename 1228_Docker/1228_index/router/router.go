package router

import (
	"fmt"

	"net/http"
	"strings"
	"test/controller"

	"github.com/gin-gonic/gin"
)

func Router() {
	r := gin.Default()
	r.GET("/livez", controller.Livez)
	r.NoRoute(gin.WrapH(http.FileServer(http.Dir("./public"))), func(c *gin.Context) {
		path := c.Request.URL.Path
		//檢查path的開頭使是否為"/"
		if !strings.HasPrefix(path, "/") {
			fmt.Println("Route Not ok")
		}
	})
	err := r.Run(":8080")
	if err != nil {
		return
	}
}
