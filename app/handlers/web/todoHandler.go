package handlers

import (
	"strconv"

	"github.com/MinadukiSekina/todo-go-app/app/domain/models"
	"github.com/MinadukiSekina/todo-go-app/app/usecases"
	"github.com/gin-gonic/gin"
)

type TodoHandler struct {
	todoUsecase usecases.TodoUsecase
}

func NewTodoHandler(uc usecases.TodoUsecase) TodoHandler {
	todoHandler := TodoHandler{todoUsecase: uc}
	return todoHandler
}

func (th *TodoHandler) Index(c *gin.Context) {
	todos, err := th.todoUsecase.Show()
	if err != nil {
		c.HTML(500, "error/error.html", gin.H{
			"message": err.Error(),
		})
		return
	}
	c.HTML(200, "todo/index.html", gin.H{
		"todos":      todos,
		"NotStarted": models.NotStarted,
		"Done":       models.Done,
	})
}

func (th *TodoHandler) ShowById(c *gin.Context) {
	id_s := c.Param("id")
	id, err := strconv.ParseUint(id_s, 10, 64)
	if err != nil {
		c.HTML(400, "error/error.html", gin.H{
			"message": err.Error(),
		})
		return
	}
	todo, err := th.todoUsecase.SearchByID(uint(id))
	if err != nil {
		c.HTML(500, "error/error.html", gin.H{
			"message": err.Error(),
		})
		return
	}
	c.HTML(200, "todo/show.html", gin.H{
		"todo":       todo,
		"NotStarted": models.NotStarted,
		"Done":       models.Done,
	})
}
