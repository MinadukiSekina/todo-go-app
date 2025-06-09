package handlers

import "github.com/gin-gonic/gin"

func SetRouting(th TodoHandler) {
	router := gin.Default()

	router.LoadHTMLGlob("templates/*/*.html")

	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "top/index.html", gin.H{"title": "Hello, World!"})
	})
	router.GET("/todo/:id", th.ShowById)
	router.GET("/todo", th.Index)
	router.POST("/todo", th.Create)
	router.Run(":3000")

}
