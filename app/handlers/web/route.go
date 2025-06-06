package handlers

import "github.com/gin-gonic/gin"

func SetRouting(th TodoHandler) {
	router := gin.Default()

	router.LoadHTMLGlob("templates/*/*.html")

	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "top/index.html", gin.H{"title": "Hello, World!"})
	})
	router.GET("/todo/:id", func(c *gin.Context) { th.SearchByID(c) })
	router.GET("/todo", func(c *gin.Context) { th.Show(c) })
	router.Run(":3000")

}
