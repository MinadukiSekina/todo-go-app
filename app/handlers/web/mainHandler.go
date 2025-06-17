package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type MainHandler struct {
}

func NewMainHandler() MainHandler {
	return MainHandler{}
}

func (mh *MainHandler) Index(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/todo")
}
