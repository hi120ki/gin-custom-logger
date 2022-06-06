package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hi120ki/gin-custom-logger/middleware"
)

func main() {
	r := gin.Default()
	r.Use(middleware.CustomLogger())

	r.GET("/string", func(c *gin.Context) {
		c.String(200, "OK")
	})

	r.GET("/json", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.Run()
}
