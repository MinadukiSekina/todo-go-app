package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ルートパスへのリクエストに対するハンドラーの構造体
type MainHandler struct {
}

// MainHandlerの新しいインスタンスを作成して返す
func NewMainHandler() MainHandler {
	return MainHandler{}
}

// ルートパスへのGetリクエストに対するハンドラー関数
func (mh *MainHandler) Index(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/todo")
}
