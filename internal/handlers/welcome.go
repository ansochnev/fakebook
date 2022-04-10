package handlers

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type WelcomePage struct {
	BasicURL    string
	pageContent []byte
}

func (h *WelcomePage) Handle(ctx *gin.Context) {
	if h.pageContent != nil {
		ctx.Data(http.StatusOK, "text/html", h.pageContent)
		return
	}

	data, err := os.ReadFile("site/welcome.html")
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Cannot read welcome page")
		return
	}

	pageContent := strings.Replace(string(data), "${base_url}", h.BasicURL, 1)
	h.pageContent = []byte(pageContent)

	h.Handle(ctx)
}
