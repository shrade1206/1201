package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
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
