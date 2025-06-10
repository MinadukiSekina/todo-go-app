package injector

import (
	"github.com/MinadukiSekina/todo-go-app/app/db"
	"github.com/MinadukiSekina/todo-go-app/app/domain/repository"
	handlers "github.com/MinadukiSekina/todo-go-app/app/handlers/web"
	"github.com/MinadukiSekina/todo-go-app/app/usecases"
)

func InjectDB() db.SqlHandler {
	sqlhandler := db.GetSqlHandler()
	return *sqlhandler
}

/*
TodoRepository(interface)に実装であるSqlHandler(struct)を渡し生成する。
*/
func InjectTodoRepository() repository.TodoRepository {
	sqlHandler := InjectDB()
	return db.NewTodoRepository(sqlHandler)
}

func InjectTodoUsecase() usecases.TodoUsecase {
	TodoRepo := InjectTodoRepository()
	return usecases.NewTodoUsecase(TodoRepo)
}

func InjectTodoHandler() handlers.TodoHandler {
	return handlers.NewTodoHandler(InjectTodoUsecase())
}
