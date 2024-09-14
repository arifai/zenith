package core

import (
	"github.com/arifai/go-modular-monolithic/internal/account/domain/model"
	"github.com/gin-gonic/gin"
)

// Context is a struct that holds the Gin context and the current account information.
type Context struct {
	Gin            *gin.Context
	CurrentAccount interface{}
}

// NewContext retrieves the account from the gin.Context and returns a new Context instance containing the account details.
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
