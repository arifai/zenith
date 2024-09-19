package core

//
//import (
//	"github.com/arifai/zenith/internal/account/domain/model"
//	"github.com/gin-gonic/gin"
//)
//
//// Context holds the Gin context and current account information.
//type Context struct {
//	Gin            *gin.Context
//	CurrentAccount interface{}
//}
//
//// NewContext returns a new Context instance with account details from the gin.Context.
//func NewContext(ctx *gin.Context) *Context {
//	account, _ := ctx.Get("account_id")
//	acc, _ := account.(*model.Account)
//
//	return &Context{
//		Gin:            ctx,
//		CurrentAccount: acc,
//	}
//}
