package route

import (
	"github.com/MinadukiSekina/todo-go-app/app/injector"
	"github.com/gin-gonic/gin"
)

func SetRouting() {
	router := gin.Default()

	router.LoadHTMLGlob("templates/*/*.html")
	router.Static("/css", "app/static/css/")

	th := injector.InjectTodoHandler()
	mh := injector.InjectMainHandler()

	router.GET("/", mh.Index)

	router.GET("/todo", th.Index)
	router.POST("/todo", th.Create)

	router.GET("/todo/:id", th.ShowById)
	router.POST("/todo/:id", th.Update)
	router.POST("/todo/:id/delete", th.Delete)

	router.Run(":3000")
}
