package main

import (
	"github.com/MinadukiSekina/todo-go-app/app/db"
	"github.com/gin-gonic/gin"
)

func main() {

	// DBの初期化
	db.Init()

	router := gin.Default()

	router.LoadHTMLGlob("templates/*/*.html")

	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "top/index.html", gin.H{"title": "Hello, World!"})
	})

	router.Run(":3000")
}
