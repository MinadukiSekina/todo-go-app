package db

import (
	"log/slog"

	"github.com/MinadukiSekina/todo-go-app/app/domain/models"
	"github.com/MinadukiSekina/todo-go-app/app/domain/repository"
)

// todoモデルのDB処理を担うリポジトリの構造体
type todoRepository struct {
	handler SqlHandler
}

// TodoRepositoryの新しいインスタンスを作成して返す
func NewTodoRepository(sqlHandler SqlHandler) repository.TodoRepository {
	todoRepository := todoRepository{handler: sqlHandler}
	return &todoRepository
}

// todoの一覧を返す
func (tr *todoRepository) FindAll() (*[]models.Todo, error) {
	var todos []models.Todo
	result := tr.handler.GetConnection().Find(&todos)
	return &todos, result.Error
}

// 指定されたIDのtodoを検索して結果を返す
func (tr *todoRepository) FindById(id uint) (*models.Todo, error) {
	var todo models.Todo
	result := tr.handler.GetConnection().Where("id = ?", id).First(&todo)
	if result.Error != nil {
		return nil, result.Error
	}
	return &todo, nil
}

// 渡されたtodoを新規作成して保存する
func (tr *todoRepository) Create(todo *models.Todo) error {
	result := tr.handler.GetConnection().Create(todo)
	return result.Error
}

// 渡されたtodoのデータを更新する
func (tr *todoRepository) Update(todo *models.Todo) error {
	// Saveメソッドだと、存在しないIDの場合はCreate動作になるため、
	// 存否チェックをする
	_, err := tr.FindById(todo.ID)
	if err != nil {
		return err
	}
	result := tr.handler.GetConnection().Save(todo)
	return result.Error
}

// 指定されたIDのtodoを削除する
func (tr *todoRepository) Delete(id uint) error {
	// 存在しないIDの場合でもエラーは出ないようなので、存否チェックをする
	_, err := tr.FindById(id)
	if err != nil {
		return err
	}
	result := tr.handler.GetConnection().Delete(&models.Todo{}, id)
	return result.Error
}

// todoRepositoryの終了処理
func (th *todoRepository) Close() error {
	// 依存先をクローズする
	err := th.handler.Close()
	if err != nil {
		slog.Error(err.Error())
	}
	return err
}
