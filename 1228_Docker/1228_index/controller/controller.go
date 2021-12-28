package controller

import "github.com/gin-gonic/gin"

func Livez(c *gin.Context) {
	c.JSON(200, nil)
}
