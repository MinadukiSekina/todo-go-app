package db

import (
	"github.com/MinadukiSekina/todo-go-app/app/domain/models"
	"github.com/MinadukiSekina/todo-go-app/app/domain/repository"
)

type todoRepository struct {
	SqlHandler
}

func NewTodoRepository(sqlHandler SqlHandler) repository.TodoRepository {
	todoRepository := todoRepository{SqlHandler: sqlHandler}
	return &todoRepository
}

func (tr *todoRepository) FindAll() (*[]models.Todo, error) {
	var todos []models.Todo
	result := tr.GetConnection().Find(&todos)
	return &todos, result.Error
}

func (tr *todoRepository) FindById(id uint) (*models.Todo, error) {
	var todo models.Todo
	result := tr.GetConnection().Where("id = ?", id).First(&todo)
	return &todo, result.Error
}

func (tr *todoRepository) Create(todo *models.Todo) error {
	result := tr.GetConnection().Create(todo)
	return result.Error
}

func (tr *todoRepository) Update(todo *models.Todo) error {
	result := tr.GetConnection().Save(todo)
	return result.Error
}

func (tr *todoRepository) Delete(id uint) error {
	result := tr.GetConnection().Delete(&models.Todo{}, id)
	return result.Error
}
