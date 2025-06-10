package usecases

import (
	"github.com/MinadukiSekina/todo-go-app/app/domain/models"
	"github.com/MinadukiSekina/todo-go-app/app/domain/repository"
)

type TodoUsecase interface {
	SearchByID(uint) (*models.Todo, error)
	Show() (todos *[]models.Todo, err error)
	Add(todo *models.Todo) error
	Edit(todo *models.Todo) error
	Delete(id uint) error
}

type todoUsecase struct {
	repos repository.TodoRepository
}

func NewTodoUsecase(todoRepo repository.TodoRepository) TodoUsecase {
	todoUsecase := todoUsecase{repos: todoRepo}
	return &todoUsecase
}

func (uc *todoUsecase) SearchByID(id uint) (todo *models.Todo, err error) {
	todo, err = uc.repos.FindById(id)
	return
}

func (uc *todoUsecase) Show() (todos *[]models.Todo, err error) {
	todos, err = uc.repos.FindAll()
	return
}

func (uc *todoUsecase) Add(todo *models.Todo) (err error) {
	err = uc.repos.Create(todo)
	return
}

func (uc *todoUsecase) Edit(todo *models.Todo) (err error) {
	err = uc.repos.Update(todo)
	return
}

func (uc *todoUsecase) Delete(id uint) (err error) {
	err = uc.repos.Delete(id)
	return
}
