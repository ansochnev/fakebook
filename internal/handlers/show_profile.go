package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"fakebook/internal/backend"
)

type ShowProfile struct {
	backend *backend.Backend
}

func (h ShowProfile) Handle(ctx *gin.Context) {
	username := ctx.Param("username")

	userProfile, err := h.backend.GetProfileByUsername(username)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	if userProfile == nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	body, _ := json.MarshalIndent(userProfile, "", "  ")
	body = append(body, byte('\n'))

	ctx.Data(http.StatusOK, "application/json", body)
}

func NewShowProfile(b *backend.Backend) gin.HandlerFunc {
	h := ShowProfile{
		backend: b,
	}
	return h.Handle
}
