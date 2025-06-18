package repository

import (
	"github.com/MinadukiSekina/todo-go-app/app/domain/interfaces"
	"github.com/MinadukiSekina/todo-go-app/app/domain/models"
)

// TodoRepository is interface for infrastructure
type TodoRepository interface {
	interfaces.Closer
	FindAll() (*[]models.Todo, error)
	FindById(id uint) (*models.Todo, error)
	Create(todo *models.Todo) error
	Update(todo *models.Todo) error
	Delete(id uint) error
}
