package repository

import (
	"github.com/arifai/zenith/internal/model"
	"github.com/arifai/zenith/pkg/common"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
	"time"
)

type (
	NotificationRepository interface {
		GetList(id *uuid.UUID, paging *common.Pagination) (notifications []*model.Notification, count int64, err error)

		MarkAsRead(id uuid.UUID) error
	}

	notificationRepository struct{ *Repository }
)

func NewNotificationRepository(r *Repository) NotificationRepository {
	return &notificationRepository{r}
}

func (r *notificationRepository) GetList(id *uuid.UUID, paging *common.Pagination) (notifications []*model.Notification, count int64, err error) {
	if err = r.db.Model(&model.Notification{}).
		Where("account_id = ?", id).
		Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err = r.db.Model(&model.Notification{}).
		Scopes(common.Paginate(paging)).
		Where("account_id = ?", id).
		Order(paging.GetSort()).
		Find(&notifications).Error; err != nil {
		return nil, 0, err
	}

	return notifications, count, nil
}

func (r *notificationRepository) MarkAsRead(id uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&model.Notification{}).
		Clauses(clause.Returning{}).
		Where(&model.Notification{ID: id}).
		Updates(map[string]interface{}{"read": true, "read_at": &now}).Error
}
