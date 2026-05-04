package service

import (
	"errors"
	"fmt"

	"github.com/example/go-microservices/notification-service/internal/models"
	"github.com/example/go-microservices/notification-service/internal/repository"
	"gorm.io/gorm"
)

type NotificationService interface {
	Create(userID uint, nType models.NotificationType, message string) (*models.Notification, error)
	ListByUser(userID uint) ([]models.Notification, error)
	GetByID(id, userID uint) (*models.Notification, error)
	MarkRead(id, userID uint) error
}

type notificationService struct {
	repo repository.NotificationRepository
}

func NewNotificationService(repo repository.NotificationRepository) NotificationService {
	return &notificationService{repo: repo}
}

func (s *notificationService) Create(userID uint, nType models.NotificationType, message string) (*models.Notification, error) {
	n := &models.Notification{UserID: userID, Type: nType, Message: message}
	if err := s.repo.Create(n); err != nil {
		return nil, fmt.Errorf("create notification: %w", err)
	}
	return n, nil
}

func (s *notificationService) ListByUser(userID uint) ([]models.Notification, error) {
	return s.repo.FindByUserID(userID)
}

func (s *notificationService) GetByID(id, userID uint) (*models.Notification, error) {
	n, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("notification not found")
		}
		return nil, err
	}
	if n.UserID != userID {
		return nil, errors.New("forbidden")
	}
	return n, nil
}

func (s *notificationService) MarkRead(id, userID uint) error {
	if _, err := s.GetByID(id, userID); err != nil {
		return err
	}
	return s.repo.MarkRead(id, userID)
}
