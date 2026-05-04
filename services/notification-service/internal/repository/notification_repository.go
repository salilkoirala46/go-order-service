package repository

import (
	"github.com/example/go-microservices/notification-service/internal/models"
	"gorm.io/gorm"
)

type NotificationRepository interface {
	Create(n *models.Notification) error
	FindByUserID(userID uint) ([]models.Notification, error)
	FindByID(id uint) (*models.Notification, error)
	MarkRead(id, userID uint) error
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(n *models.Notification) error {
	return r.db.Create(n).Error
}

func (r *notificationRepository) FindByUserID(userID uint) ([]models.Notification, error) {
	var notifications []models.Notification
	if err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&notifications).Error; err != nil {
		return nil, err
	}
	return notifications, nil
}

func (r *notificationRepository) FindByID(id uint) (*models.Notification, error) {
	var n models.Notification
	if err := r.db.First(&n, id).Error; err != nil {
		return nil, err
	}
	return &n, nil
}

func (r *notificationRepository) MarkRead(id, userID uint) error {
	return r.db.Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("read", true).Error
}
