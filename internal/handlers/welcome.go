package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type WelcomePage struct {
}

func (h WelcomePage) Handle(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Welcome to Fakebook\n")
}

func NewWelcomePage() gin.HandlerFunc {
	wp := WelcomePage{}
	return wp.Handle
}
