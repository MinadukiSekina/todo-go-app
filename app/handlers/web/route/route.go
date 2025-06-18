package route

import (
	"github.com/MinadukiSekina/todo-go-app/app/injector"
	"github.com/gin-gonic/gin"
)

// ルーティングやミドルウェアの設定を行う
func SetRouting() {
	router := gin.Default()

	// HTML・css・jsファイルの読み込み
	router.LoadHTMLGlob("app/templates/*/*.html")
	router.Static("/css", "app/static/css/")
	router.Static("/js", "app/static/js/")

	// 各ハンドラーの初期化
	th := injector.InjectTodoHandler()
	mh := injector.InjectMainHandler()

	// ハンドラーの終了処理
	defer th.Close()

	// ルーティングの設定
	router.GET("/", mh.Index)

	router.GET("/todo", th.Index)
	router.POST("/todo", th.Create)

	router.GET("/todo/:id", th.ShowById)
	router.POST("/todo/:id", th.Update)
	router.POST("/todo/:id/delete", th.Delete)

	// 待機開始
	router.Run(":3000")
}
