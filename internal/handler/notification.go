package handler

import (
	"github.com/arifai/zenith/internal/service"
	"github.com/arifai/zenith/internal/types/request"
	"github.com/arifai/zenith/pkg/common"
	"github.com/arifai/zenith/pkg/utils"
	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	*Handler
	notificationService service.NotificationService
}

func NewNotificationHandler(handler *Handler, notificationService service.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		Handler:             handler,
		notificationService: notificationService,
	}
}

func (h *NotificationHandler) GetList(ctx *gin.Context) {
	paging, err := utils.ValidateQuery[common.Pagination](ctx)
	if err != nil {
		h.response.Error(ctx, err)
		return
	}

	accountID := GetAccountIDFromContext(ctx)
	if accountID == nil {
		h.response.NotFound(ctx, "Account ID not found in context")
		return
	}

	entries, err := h.notificationService.GetList(accountID, paging)
	if err != nil {
		h.response.Error(ctx, err)
		return
	}

	h.response.Success(ctx, entries)
}

func (h *NotificationHandler) MarkAsRead(ctx *gin.Context) {
	body, err := utils.ValidateBody[request.NotificationMarkAsReadRequest](ctx)
	if err != nil {
		h.response.Error(ctx, err)
		return
	}

	founded, err := h.notificationService.MarkAsRead(body.ID)
	if err != nil {
		h.response.Error(ctx, err)
		return
	}

	if !founded {
		h.response.NotFound(ctx, "notification not found")
		return
	}

	h.response.Success(ctx, nil)
}
