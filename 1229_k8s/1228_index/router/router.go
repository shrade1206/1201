package router

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Router() {
	r := gin.Default()
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		ps := strings.Split(path, "/")
		fs := gin.Dir("./public", true)
		_, err := fs.Open(ps[1])
		if err != nil {
			c.File("./public/index.html")
		} else {
			h := http.FileServer(http.Dir("./public"))
			h.ServeHTTP(c.Writer, c.Request)
		}
	})
	err := r.Run(":8080")
	if err != nil {
		return
	}
}
