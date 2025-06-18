package usecases

import (
	"log/slog"

	"github.com/MinadukiSekina/todo-go-app/app/domain/interfaces"
	"github.com/MinadukiSekina/todo-go-app/app/domain/models"
	"github.com/MinadukiSekina/todo-go-app/app/domain/repository"
)

// ユースケースのインターフェイス
type TodoUsecase interface {
	interfaces.Closer
	SearchByID(uint) (*models.Todo, error)
	Show() (todos *[]models.Todo, err error)
	Add(todo *models.Todo) error
	Edit(todo *models.Todo) error
	Delete(id uint) error
}

// todoに関わるユースケースの構造体
type todoUsecase struct {
	repos repository.TodoRepository
}

// TodoUsecaseの新しいインスタンスを作成して返す
func NewTodoUsecase(todoRepo repository.TodoRepository) TodoUsecase {
	todoUsecase := todoUsecase{repos: todoRepo}
	return &todoUsecase
}

// 指定されたIDのtodoを検索して結果を返す
func (uc *todoUsecase) SearchByID(id uint) (todo *models.Todo, err error) {
	todo, err = uc.repos.FindById(id)
	return
}

// todoの一覧を検索して返す
func (uc *todoUsecase) Show() (todos *[]models.Todo, err error) {
	todos, err = uc.repos.FindAll()
	return
}

// 渡されたtodoを新規作成して保存する
func (uc *todoUsecase) Add(todo *models.Todo) (err error) {
	err = uc.repos.Create(todo)
	return
}

// 渡されたtodoを更新して保存する
func (uc *todoUsecase) Edit(todo *models.Todo) (err error) {
	err = uc.repos.Update(todo)
	return
}

// 指定されたIDのtodoを削除する
func (uc *todoUsecase) Delete(id uint) (err error) {
	err = uc.repos.Delete(id)
	return
}

// ユースケースの終了処理を行う
func (uc *todoUsecase) Close() error {
	err := uc.repos.Close()
	if err != nil {
		slog.Error(err.Error())
	}
	return err
}
