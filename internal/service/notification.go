package service

import (
	"github.com/arifai/zenith/internal/model"
	"github.com/arifai/zenith/internal/repository"
	"github.com/arifai/zenith/pkg/common"
	"github.com/google/uuid"
)

type (
	NotificationService interface {
		GetList(id *uuid.UUID, paging *common.Pagination) (*common.EntriesModel[*model.Notification], error)

		MarkAsRead(id uuid.UUID) error
	}

	notificationService struct {
		*Service
		notificationRepo repository.NotificationRepository
	}
)

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

func (s *notificationService) MarkAsRead(id uuid.UUID) error {
	return s.notificationRepo.MarkAsRead(id)
}
