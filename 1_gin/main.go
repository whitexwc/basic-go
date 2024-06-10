package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.POST("/post", func(c *gin.Context) {
		c.String(http.StatusOK, "hello post method")
	})

	router.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "parameter router "+name)

	})

	router.GET("/views/*.html", func(c *gin.Context) {
		page := c.Param(".html")
		c.String(http.StatusOK, "hello 通配符路由 "+page)
	})

	router.GET("/order", func(c *gin.Context) {
		id := c.Query("id")
		c.String(http.StatusOK, "id is:"+id)
	})

	router.Run()
}
