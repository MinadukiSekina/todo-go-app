package handlers

import (
	"log/slog"
	"net/http"
	"sort"
	"strconv"

	"github.com/MinadukiSekina/todo-go-app/app/domain/models"
	"github.com/MinadukiSekina/todo-go-app/app/usecases"
	"github.com/gin-gonic/gin"
)

// todoパスへのリクエストに対するハンドラーの構造体
type TodoHandler struct {
	todoUsecase usecases.TodoUsecase
}

// TodoHandlerの新しいインスタンスを作成して返す
func NewTodoHandler(uc usecases.TodoUsecase) TodoHandler {
	todoHandler := TodoHandler{todoUsecase: uc}
	return todoHandler
}

// todoの一覧を表示する
func (th *TodoHandler) Index(c *gin.Context) {
	todos, err := th.todoUsecase.Show()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error/error.html", gin.H{
			"message": err.Error(),
		})
		return
	}

	// Sort todos by status in ascending order
	sort.Slice(*todos, func(i, j int) bool {
		return (*todos)[i].Status < (*todos)[j].Status
	})

	// Get flash message content if it exists
	fm := GetFlashMessage(c)

	c.HTML(http.StatusOK, "todo/index.html", gin.H{
		"todos":      todos,
		"NotStarted": models.NotStarted,
		"Done":       models.Done,
		flashMessage: fm.Message,
		flashType:    fm.Type,
	})
}

// 指定されたIDのtodoを表示する
func (th *TodoHandler) ShowById(c *gin.Context) {
	id_s := c.Param("id")
	id, err := strconv.ParseUint(id_s, 10, 64)
	if err != nil {
		SetFlashMessage(c, resultIsError, "このタスクは閲覧できません。")
		c.Redirect(http.StatusSeeOther, "/todo")
		return
	}
	todo, err := th.todoUsecase.SearchByID(uint(id))
	if err != nil {
		SetFlashMessage(c, resultIsError, "該当するタスクが見つかりませんでした。")
		c.Redirect(http.StatusSeeOther, "/todo")
		return
	}

	// Get flash message content if it exists
	fm := GetFlashMessage(c)

	c.HTML(http.StatusOK, "todo/show.html", gin.H{
		"todo":       todo,
		"NotStarted": models.NotStarted,
		"Done":       models.Done,
		flashMessage: fm.Message,
		flashType:    fm.Type,
	})
}

// todoを新規作成する
func (th *TodoHandler) Create(c *gin.Context) {
	title := c.PostForm("title")
	todo := models.Todo{Title: title, Status: models.NotStarted}
	err := th.todoUsecase.Add(&todo)
	if err != nil {
		SetFlashMessage(c, resultIsError, "新しいタスクの作成に失敗しました。")
		c.Redirect(http.StatusSeeOther, "/todo")
		return
	}
	SetFlashMessage(c, resultIsSuccess, "新しいタスクを作成しました。")
	c.Redirect(http.StatusFound, "/todo")
}

func (th *TodoHandler) Update(c *gin.Context) {
	id_s := c.Param("id")
	id, err := strconv.ParseUint(id_s, 10, 64)
	if err != nil {
		SetFlashMessage(c, resultIsError, "このタスクは更新できません。")
		c.Redirect(http.StatusSeeOther, "/todo")
		return
	}
	title := c.PostForm("title")
	status_s := c.PostForm("status")

	// 文字列のstatusをStatus型に変換
	var correspond = map[string]models.Status{
		"notStarted": models.NotStarted,
		"completed":  models.Done,
	}
	status, err := models.StrToStatus(status_s, correspond)
	if err != nil {
		SetFlashMessage(c, resultIsError, "タスクの状態が不正な値です。")
		c.Redirect(http.StatusSeeOther, "/todo/"+id_s)
		return
	}

	// 既存のTodoを取得
	existingTodo, err := th.todoUsecase.SearchByID(uint(id))
	if err != nil {
		SetFlashMessage(c, resultIsError, "対象となるタスクが存在しません。")
		c.Redirect(http.StatusSeeOther, "/todo")
		return
	}

	existingTodo.Title = title
	existingTodo.Status = status
	err = th.todoUsecase.Edit(existingTodo)
	if err != nil {
		SetFlashMessage(c, resultIsError, "タスクの内容を更新できませんでした。")
		c.Redirect(http.StatusSeeOther, "/todo/"+id_s)
		return
	}
	SetFlashMessage(c, resultIsSuccess, "タスクの内容を更新しました。")
	c.Redirect(http.StatusFound, "/todo/"+id_s)
}

// todoを削除する
func (th *TodoHandler) Delete(c *gin.Context) {
	id_s := c.Param("id")
	id, err := strconv.ParseUint(id_s, 10, 64)
	if err != nil {
		SetFlashMessage(c, resultIsError, "このタスクは削除できません。")
		c.Redirect(http.StatusSeeOther, "/todo")
		return
	}

	err = th.todoUsecase.Delete(uint(id))

	if err != nil {
		SetFlashMessage(c, resultIsError, "削除できませんでした。")
		c.Redirect(http.StatusSeeOther, "/todo/"+id_s)
		return
	}
	SetFlashMessage(c, resultIsSuccess, "タスクの削除を完了しました。")
	c.Redirect(http.StatusFound, "/todo")
}

// 終了処理を行う
func (th *TodoHandler) Close() {
	err := th.todoUsecase.Close()
	if err != nil {
		slog.Error(err.Error())
	}
}
