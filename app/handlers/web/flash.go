package handlers

import (
	"github.com/gin-gonic/gin"
)

type FlashMessage struct {
	Type    string
	Message string
}

func SetFlashMessage(c *gin.Context, status string, msg string) {
	c.SetCookie(flashType, status, 1, "/", "localhost", true, true)
	c.SetCookie(flashMessage, msg, 1, "/", "localhost", true, true)
}

func GetFlashMessage(c *gin.Context) FlashMessage {
	t, _ := c.Cookie(flashType)
	m, _ := c.Cookie(flashMessage)
	return FlashMessage{Type: t, Message: m}
}

const (
	resultIsSuccess string = "success"
	resultIsError   string = "error"
)

const (
	flashType    string = "Type"
	flashMessage string = "Message"
)
