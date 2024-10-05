package repository

import (
	"github.com/arifai/zenith/internal/model"
	"github.com/arifai/zenith/pkg/common"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
	"time"
)

type (
	// NotificationRepository defines methods for interacting with notifications in a data store.
	// GetList fetches a list of notifications and the total count based on account ID and pagination info.
	// MarkAsRead marks a specific notification as read by its ID and returns if it was found and updated successfully.
	NotificationRepository interface {

		// GetList fetches a list of notifications and the total count for a given account ID and pagination parameters.
		GetList(id *uuid.UUID, paging *common.Pagination) (notifications []*model.Notification, count int64, err error)

		// MarkAsRead marks a specific notification as read by its ID and returns if it was found and updated successfully.
		MarkAsRead(id string) (founded bool, err error)
	}

	// notificationRepository implements NotificationRepository interface, provides repository functions for notifications.
	notificationRepository struct{ *Repository }
)

// NewNotificationRepository creates a new instance of NotificationRepository with the provided Repository parameter.
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

func (r *notificationRepository) MarkAsRead(id string) (founded bool, err error) {
	now := time.Now()
	parseID, err := uuid.Parse(id)
	if err != nil {
		return false, err
	}

	result := r.db.Model(&model.Notification{}).
		Clauses(clause.Returning{}).
		Where(&model.Notification{ID: parseID}).
		Updates(map[string]interface{}{"read": true, "read_at": &now})

	if result.Error != nil {
		return false, result.Error
	}

	if result.RowsAffected == 0 {
		return false, err
	}

	return true, nil
}
