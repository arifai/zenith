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

		MarkAsRead(id uuid.UUID) (founded bool, err error)
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
		Scopes(common.Paginate(paging, "title")).
		Where("account_id = ?", id).
		Find(&notifications).Error; err != nil {
		return nil, 0, err
	}

	count = int64(len(notifications))

	return notifications, count, nil
}

func (r *notificationRepository) MarkAsRead(id uuid.UUID) (founded bool, err error) {
	now := time.Now()
	result := r.db.Model(&model.Notification{}).
		Clauses(clause.Returning{}).
		Where(&model.Notification{ID: id}).
		Updates(map[string]interface{}{"read": true, "read_at": &now})

	if result.Error != nil {
		return false, result.Error
	}

	if result.RowsAffected == 0 {
		return false, err
	}

	return true, nil
}
