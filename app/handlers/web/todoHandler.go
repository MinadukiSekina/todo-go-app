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

func (th *TodoHandler) SearchByID(c *gin.Context) {
	n := c.Param("id")
	id, err := strconv.ParseUint(n, 10, 64)
	if err != nil {
		panic(err)
	}
	todo, err := th.todoUsecase.SearchByID(uint(id))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, todo)
}

func (th *TodoHandler) Show(c *gin.Context) {
	todos, err := th.todoUsecase.Show()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.HTML(200, "todo/show.html", gin.H{
		"todos":      todos,
		"NotStarted": models.NotStarted,
		"Done":       models.Done,
	})
}
