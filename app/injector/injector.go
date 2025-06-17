package injector

import (
	"github.com/MinadukiSekina/todo-go-app/app/db"
	"github.com/MinadukiSekina/todo-go-app/app/domain/repository"
	handlers "github.com/MinadukiSekina/todo-go-app/app/handlers/web"
	"github.com/MinadukiSekina/todo-go-app/app/usecases"
)

// データベース接続を初期化し、SqlHandlerを返す
func InjectDB() db.SqlHandler {
	sqlhandler := db.GetSqlHandler()
	return sqlhandler
}

// sqlHandlerを使用してTodoRepositoryを生成する
func InjectTodoRepository() repository.TodoRepository {
	sqlHandler := InjectDB()
	return db.NewTodoRepository(sqlHandler)
}

// TodoRepositoryを使用してTodoUsecaseを生成する
func InjectTodoUsecase() usecases.TodoUsecase {
	TodoRepo := InjectTodoRepository()
	return usecases.NewTodoUsecase(TodoRepo)
}

// TodoUsecaseを使用してTodoHandlerを生成する
func InjectTodoHandler() handlers.TodoHandler {
	return handlers.NewTodoHandler(InjectTodoUsecase())
}

// アプリケーションのメインハンドラー（ルートパス用）を生成する
func InjectMainHandler() handlers.MainHandler {
	return handlers.NewMainHandler()
}
