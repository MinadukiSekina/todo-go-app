package db

import (
	"github.com/MinadukiSekina/todo-go-app/app/domain/models"
	"github.com/MinadukiSekina/todo-go-app/app/domain/repository"
)

type TodoRepository struct {
	SqlHandler
}

func NewTodoRepository(sqlHandler SqlHandler) repository.TodoRepository {
	todoRepository := TodoRepository{sqlHandler}
	return &todoRepository
}

func (tr *TodoRepository) FindAll() (*[]models.Todo, error) {
	var todos []models.Todo
	result := tr.GetConnection().Find(&todos)
	return &todos, result.Error
}

func (tr *TodoRepository) FindById(id uint) (*models.Todo, error) {
	var todo models.Todo
	result := tr.GetConnection().Where("id = ?", id).Find(&todo)
	return &todo, result.Error
}

func (tr *TodoRepository) Create(todo *models.Todo) error {
	result := tr.GetConnection().Create(todo)
	return result.Error
}

func (tr *TodoRepository) Update(todo *models.Todo) error {
	result := tr.GetConnection().Save(todo)
	return result.Error
}

func (tr *TodoRepository) Delete(todo *models.Todo) error {
	result := tr.GetConnection().Delete(todo)
	return result.Error
}
