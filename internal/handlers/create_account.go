package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"fakebook/internal/backend"
)

type CreateAccount struct {
	Backend *backend.Backend
}

func (h *CreateAccount) Handle(ctx *gin.Context) {
	createAccountREQ, err := h.newCreateAccountREQ(ctx)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	account, err := h.Backend.CreateAccount(createAccountREQ)
	if err != nil {
		h.errorResponse(ctx, err)
		return
	}

	userHome := fmt.Sprintf("http://localhost:8080/%s", account.Username)
	ctx.Redirect(http.StatusSeeOther, userHome)
}

func (h *CreateAccount) newCreateAccountREQ(ctx *gin.Context) (*backend.CreateAccountREQ, error) {
	createAccountREQ := &backend.CreateAccountREQ{}
	var exists bool

	createAccountREQ.Email, exists = ctx.GetPostForm("email")
	if !exists {
		return nil, backend.MissingParamError("email")
	}

	createAccountREQ.FirstName, exists = ctx.GetPostForm("firstName")
	if !exists {
		return nil, backend.MissingParamError("firstName")
	}

	createAccountREQ.LastName, exists = ctx.GetPostForm("lastName")
	if !exists {
		return nil, backend.MissingParamError("lastName")
	}

	createAccountREQ.Password, exists = ctx.GetPostForm("password")
	if !exists {
		return nil, backend.MissingParamError("password")
	}

	createAccountREQ.Password2, exists = ctx.GetPostForm("password2")
	if !exists {
		return nil, backend.MissingParamError("password2")
	}

	createAccountREQ.Username = ctx.PostForm("username")

	return createAccountREQ, nil
}

func (h *CreateAccount) errorResponse(ctx *gin.Context, e error) {
	err, ok := e.(*backend.Error)
	if !ok {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	switch err.Code() {
	case backend.ErrInternal:
		ctx.JSON(http.StatusInternalServerError, err)
	default:
		ctx.JSON(http.StatusBadRequest, err)
	}
}
