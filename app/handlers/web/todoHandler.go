package handlers

import (
	"net/http"
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
		c.HTML(http.StatusInternalServerError, "error/error.html", gin.H{
			"message": err.Error(),
		})
		return
	}
	c.HTML(http.StatusOK, "todo/index.html", gin.H{
		"todos":      todos,
		"NotStarted": models.NotStarted,
		"Done":       models.Done,
	})
}

func (th *TodoHandler) ShowById(c *gin.Context) {
	id_s := c.Param("id")
	id, err := strconv.ParseUint(id_s, 10, 64)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error/error.html", gin.H{
			"message": err.Error(),
		})
		return
	}
	todo, err := th.todoUsecase.SearchByID(uint(id))
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error/error.html", gin.H{
			"message": err.Error(),
		})
		return
	}
	c.HTML(http.StatusOK, "todo/show.html", gin.H{
		"todo":       todo,
		"NotStarted": models.NotStarted,
		"Done":       models.Done,
	})
}

func (th *TodoHandler) Create(c *gin.Context) {
	title := c.PostForm("title")
	todo := models.Todo{Title: title}
	err := th.todoUsecase.Add(&todo)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error/error.html", gin.H{
			"message": err.Error(),
		})
		return
	}
	c.Redirect(http.StatusFound, "/todo")
}
