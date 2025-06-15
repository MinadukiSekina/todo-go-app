package db

import (
	"github.com/MinadukiSekina/todo-go-app/app/domain/models"
	"github.com/MinadukiSekina/todo-go-app/app/domain/repository"
)

type todoRepository struct {
	handler SqlHandler
}

func NewTodoRepository(sqlHandler SqlHandler) repository.TodoRepository {
	todoRepository := todoRepository{handler: sqlHandler}
	return &todoRepository
}

func (tr *todoRepository) FindAll() (*[]models.Todo, error) {
	var todos []models.Todo
	result := tr.handler.GetConnection().Find(&todos)
	return &todos, result.Error
}

func (tr *todoRepository) FindById(id uint) (*models.Todo, error) {
	var todo models.Todo
	result := tr.handler.GetConnection().Where("id = ?", id).First(&todo)
	if result.Error != nil {
		return nil, result.Error
	}
	return &todo, nil
}

func (tr *todoRepository) Create(todo *models.Todo) error {
	result := tr.handler.GetConnection().Create(todo)
	return result.Error
}

func (tr *todoRepository) Update(todo *models.Todo) error {
	result := tr.handler.GetConnection().Save(todo)
	return result.Error
}

func (tr *todoRepository) Delete(id uint) error {
	result := tr.handler.GetConnection().Delete(&models.Todo{}, id)
	return result.Error
}
