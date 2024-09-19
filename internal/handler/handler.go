package handler

import (
	"github.com/arifai/zenith/pkg/common"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler is a struct used to handle HTTP responses using common.Response.
type Handler struct {
	response *common.Response
}

// New initializes and returns a new instance of Handler by embedding the provided Response.
func New(response *common.Response) *Handler {
	return &Handler{response: response}
}

// GetAccountIDFromContext retrieves the account ID from the provided gin.Context.
// If the account ID does not exist in the context, an empty string is returned.
func GetAccountIDFromContext(ctx *gin.Context) *uuid.UUID {
	id, exists := ctx.Get("account_id")
	if !exists {
		return nil
	}

	accountId, ok := id.(*uuid.UUID)
	if !ok {
		return nil
	}

	return accountId
}
