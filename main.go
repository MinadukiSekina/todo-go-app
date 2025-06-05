package main

import "github.com/gin-gonic/gin"

func main() {
	router := gin.Default()

	router.LoadHTMLGlob("templates/*/*.html")

	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "top/index.html", gin.H{"title": "Hello, World!"})
	})

	router.Run(":3000")
}
