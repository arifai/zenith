package service

import (
	"github.com/arifai/zenith/internal/model"
	"github.com/arifai/zenith/internal/repository"
	"github.com/arifai/zenith/pkg/common"
	"github.com/arifai/zenith/pkg/errormessage"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type (
	// NotificationService provides methods for managing notifications.
	// GetList retrieves a list of notifications and their pagination metadata based on the given ID and pagination parameters.
	// MarkAsRead marks a notification as read by its ID, indicating if the operation was successful and if the notification was found.
	NotificationService interface {
		// GetList retrieves a list of notifications and pagination details based on the given account ID and pagination parameters.
		GetList(id *uuid.UUID, paging *common.Pagination) (*common.EntriesModel[*model.Notification], error)

		// MarkAsRead marks a notification as read by its ID, returning if it was found and any error encountered.
		MarkAsRead(id string) (founded bool, err error)
	}

	// notificationService struct implements the NotificationService interface, providing methods to manage notifications.
	notificationService struct {
		*Service
		notificationRepo repository.NotificationRepository
	}
)

// NewNotificationService creates a new instance of NotificationService with the provided service and NotificationRepository.
func NewNotificationService(service *Service, notificationRepo repository.NotificationRepository) NotificationService {
	return &notificationService{
		Service:          service,
		notificationRepo: notificationRepo,
	}
}

func (s *notificationService) GetList(id *uuid.UUID, paging *common.Pagination) (*common.EntriesModel[*model.Notification], error) {
	entries, count, err := s.notificationRepo.GetList(id, paging)
	if err != nil {
		return nil, err
	}

	page := paging.GetPage(count)
	totalPages := paging.GetTotalPages(count)

	return common.NewEntries(entries, count, page, totalPages), nil
}

func (s *notificationService) MarkAsRead(id string) (founded bool, err error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		s.log.Error(errormessage.ErrFailedToParseUUIDText, zap.String("input", id), zap.Error(err))
		return
	}
	return s.notificationRepo.MarkAsRead(parsedID)
}
