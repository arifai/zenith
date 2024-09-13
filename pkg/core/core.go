package core

import (
	"github.com/arifai/go-modular-monolithic/internal/account/domain/model"
	"github.com/gin-gonic/gin"
)

// Context is a struct to define the context
type Context struct {
	Gin            *gin.Context
	CurrentAccount interface{}
}

// NewContext is a function to create a new context
func NewContext(ctx *gin.Context) *Context {
	account, exists := ctx.Get("account")
	if !exists {
		return &Context{
			Gin:            ctx,
			CurrentAccount: nil,
		}
	}

	acc, ok := account.(*model.Account)
	if !ok {
		return &Context{
			Gin:            ctx,
			CurrentAccount: nil,
		}
	}

	return &Context{
		Gin:            ctx,
		CurrentAccount: acc,
	}
}
