package handlers

import (
	"github.com/gin-gonic/gin"
)

// フラッシュメッセージの内容を示す構造体
type FlashMessage struct {
	Type    string
	Message string
}

// cookieにフラッシュメッセージの内容を保存する
func SetFlashMessage(c *gin.Context, status string, msg string) {
	c.SetCookie(flashType, status, 1, "/", "localhost", true, true)
	c.SetCookie(flashMessage, msg, 1, "/", "localhost", true, true)
}

// cookieからフラッシュメッセージの内容を取り出して構造体に変換する
func GetFlashMessage(c *gin.Context) FlashMessage {
	t, _ := c.Cookie(flashType)
	m, _ := c.Cookie(flashMessage)
	return FlashMessage{Type: t, Message: m}
}

// フラッシュメッセージの種別
const (
	resultIsSuccess string = "success"
	resultIsError   string = "error"
)

// cookieに設定する値の名前
const (
	flashType    string = "Type"
	flashMessage string = "Message"
)
