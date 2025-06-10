package main

import (
	"github.com/MinadukiSekina/todo-go-app/app/db"
	handlers "github.com/MinadukiSekina/todo-go-app/app/handlers/web"
	"github.com/MinadukiSekina/todo-go-app/app/injector"
)

func main() {

	// DBの初期化
	db.Init()
	todoHandler := injector.InjectTodoHandler()
	handlers.SetRouting(todoHandler)
}
